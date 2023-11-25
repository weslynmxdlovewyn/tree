package template

import (
	"errors"
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/rolandhe/daog"
	txrequest "github.com/rolandhe/daog/tx"
	"github.com/weslynmxdlovewyn/tree/repository/tpl"
)

var noActiveErr = errors.New("no active template")

// getActiveTpl 从数据库里查询出正在使用的模板
func getActiveTpl(dc *dgctx.DgContext) (*tpl.QuestionsTpl, error) {
	return daog.AutoTransWithResult[*tpl.QuestionsTpl](func() (*daog.TransContext, error) {
		return daog.NewTransContext(nil, txrequest.RequestReadonly, dc.TraceId)
	}, func(tc *daog.TransContext) (*tpl.QuestionsTpl, error) {
		all, err := tpl.QuestionsTplDao.GetAll(tc, tpl.QuestionsTplFields.Id, tpl.QuestionsTplFields.IsSealed)
		if err != nil {
			return nil, err
		}
		var ids []int64
		for _, tpl := range all {
			if tpl.IsSealed == 0 {
				ids = append(ids, tpl.Id)
			}
		}
		if len(ids) == 0 {
			return nil, noActiveErr
		}
		if len(ids) > 1 {
			return nil, errors.New("need only on active template")
		}
		return tpl.QuestionsTplDao.GetById(tc, ids[0])
	})
}

func getTplById(dc *dgctx.DgContext, id int64) (*tpl.QuestionsTpl, error) {
	return daog.AutoTransWithResult[*tpl.QuestionsTpl](func() (*daog.TransContext, error) {
		return daog.NewTransContext(nil, txrequest.RequestReadonly, dc.TraceId)
	}, func(tc *daog.TransContext) (*tpl.QuestionsTpl, error) {
		return tpl.QuestionsTplDao.GetById(tc, id)
	})
}
