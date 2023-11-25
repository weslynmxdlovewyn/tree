package template

import (
	"encoding/json"
	"errors"
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/weslynmxdlovewyn/tree/dy_schema"
	"github.com/weslynmxdlovewyn/tree/repository/tpl"
	"github.com/weslynmxdlovewyn/tree/utils"
	"sync"
	"time"
)

const cacheTime = 300

var cache = &dySchemaCache{
	sealedTplCache: sync.Map{},
	activeTpl: &activeLoader{
		expiredAt: 0,
		commonLoader: commonLoader{
			loadFunc: getActiveFromDB,
		},
	},
}

type dySchemaCache struct {
	sealedTplCache sync.Map
	activeTpl      *activeLoader
}

type sealedLoader struct {
	commonLoader
	err  error
	wait chan struct{}
}

func (s *sealedLoader) get(dc *dgctx.DgContext) (*DbTemplate, error) {
	<-s.wait
	return s.dbTpl, s.err
}

func (dyc *dySchemaCache) getSealedTplById(dc *dgctx.DgContext, id int64) (*DbTemplate, error) {
	load, ok := dyc.sealedTplCache.Load(id)
	if ok {
		tl := load.(*sealedLoader)
		return tl.get(dc)
	}

	newLoader := &sealedLoader{
		wait: make(chan struct{}),
		commonLoader: commonLoader{
			loadFunc: func(dc *dgctx.DgContext, dbTpl *DbTemplate) error {
				return getSealedFromDb(dc, id, dbTpl)
			},
		},
	}

	old, loaded := dyc.sealedTplCache.LoadOrStore(id, newLoader)
	if loaded {
		return old.(*sealedLoader).get(dc)
	}

	err := newLoader.load(dc)
	newLoader.err = err
	close(newLoader.wait)
	if err != nil {
		dyc.sealedTplCache.Delete(id)
		return nil, err
	}
	return newLoader.dbTpl, err
}

type activeLoader struct {
	sync.RWMutex
	commonLoader
	expiredAt int64
}

func (a *activeLoader) get(dc *dgctx.DgContext) (*DbTemplate, error) {
	a.RLock()
	if time.Now().Unix() <= a.expiredAt {
		defer a.RUnlock()
		return a.dbTpl, nil
	}

	a.RUnlock()

	a.Lock()
	defer a.Unlock()
	if time.Now().Unix() <= a.expiredAt {
		return a.dbTpl, nil
	}

	err := a.load(dc)
	if err != nil {
		return nil, err
	}

	a.expiredAt = time.Now().Unix() + cacheTime

	return a.dbTpl, nil
}

type commonLoader struct {
	dbTpl    *DbTemplate
	loadFunc func(dc *dgctx.DgContext, dbTpl *DbTemplate) error
}

func (c *commonLoader) load(dc *dgctx.DgContext) error {
	var err error
	dbTpl := &DbTemplate{}
	err = c.loadFunc(dc, dbTpl)
	if err != nil {
		return err
	}
	// 填充  LeafInfosMap  SimpleLeafList 两个字段
	dbTpl.Supplement()

	c.dbTpl = dbTpl
	return nil
}

type DbTemplate struct {
	Id              int64
	Name            string
	QuestionsSchema *dy_schema.QuestionsDySchema[dy_schema.QuestionDesc] // 对应表字段text
	Layout          string                                               // 什么意思 作用
	LeafInfosMap    map[string]*dy_schema.LeafInfo                       // 业务需要，对应表的字段page_layout
	SimpleLeafList  []*dy_schema.CategoryBasic                           // 什么意思 作用
	SchemaSha1      string                                               // 什么意思 作用
	AllKeywords     *TplKeywords                                         // 业务需要，对应表的字段keywords
}

// Supplement 这个方法的作用是为了填充其他字段
func (dt *DbTemplate) Supplement() {
	leafs := dt.QuestionsSchema.GetLeafs()

	dt.SimpleLeafList = make([]*dy_schema.CategoryBasic, 0, len(leafs))

	mp := make(map[string]*dy_schema.LeafInfo, len(leafs))

	for _, info := range leafs {
		mp[info.Code] = info
		dt.SimpleLeafList = append(dt.SimpleLeafList, info.CategoryBasic)
	}
	dt.LeafInfosMap = mp
}

type TplKeywords struct {
	MatchKeywords []string `json:"matchKeywords"`
	AiKeywords    []string `json:"aiKeywords"`
	ToJobKeys     []string `json:"toJobKeys"`
	AiGenJdKeys   []string `json:"aiGenJdKeys"`
}

func GetByTemplateId(dc *dgctx.DgContext, id int64) (*DbTemplate, error) {
	info, err := cache.activeTpl.get(dc)
	if err != nil {
		return nil, err
	}
	if info.Id == id {
		return info, nil
	}

	return cache.getSealedTplById(dc, id)
}

// 从数据库中查询模板内容，并填充到字段 dbTpl（dbTpl事先不能是nil）
func getActiveFromDB(dc *dgctx.DgContext, dbTpl *DbTemplate) error {
	tplQuestion, err := getActiveTpl(dc)
	if err != nil {
		return err
	}

	return fillDbTemplate(dbTpl, tplQuestion)
}

func getSealedFromDb(dc *dgctx.DgContext, id int64, dbtpl *DbTemplate) error {
	tplQuestion, err := getTplById(dc, id)
	if err != nil {
		return err
	}
	if tplQuestion == nil {
		return errors.New("no specified template")
	}
	if tplQuestion.IsSealed != 1 {
		return errors.New("template in used")
	}

	return fillDbTemplate(dbtpl, tplQuestion)
}

// 用 tplQuestion的内容 填充 dbTpl
func fillDbTemplate(dbTpl *DbTemplate, tplQuestion *tpl.QuestionsTpl) error {
	schema := &dy_schema.QuestionsDySchema[dy_schema.QuestionDesc]{}
	binText := []byte(tplQuestion.Text)
	_ = json.Unmarshal(binText, schema)
	dbTpl.SchemaSha1 = utils.Sha1Data(binText)

	if tplQuestion.Keywords.Valid {
		tplKeywords := &TplKeywords{}
		_ = json.Unmarshal([]byte(tplQuestion.Keywords.String), tplKeywords)
		dbTpl.AllKeywords = tplKeywords
	}

	dbTpl.Id = tplQuestion.Id
	dbTpl.Name = tplQuestion.Name
	dbTpl.QuestionsSchema = schema
	dbTpl.Layout = tplQuestion.PageLayout.StringNilAsEmpty()

	return nil
}
