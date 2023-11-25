package dy_schema

import (
	"github.com/weslynmxdlovewyn/tree/utils"
	"strings"
)

const (
	walkQuestion = 1
	walkLeaf     = 2
	walkAll      = 3
)

type QuestionsDySchema[T any] struct {
	Leaf          int
	MaxQuestionId int
	Contents      []*Category[T]
}

type Category[T any] struct {
	CategoryBasic
	Children   []*Category[T] `json:"children,omitempty"`
	Questions  []*T           `json:"questions,omitempty"`
	Appendable bool           `json:"appendable"`
	Replace    bool           `json:"replace"`
}

type CategoryBasic struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// question ....

type QuestionDesc struct {
	Id        int    `json:"id"`
	Desc      string `json:"desc"`
	IsDefault bool   `json:"isDefault"`
	ParenCode string `json:"parentCode,omitempty"`
}

// leafInfo

type LeafInfo struct {
	Path string `json:"path"`
	*CategoryBasic
	Replace bool
}

func (qs *QuestionsDySchema[T]) GetLeafs() []*LeafInfo {
	ctx := walkCtx[T]{
		pathStack: &utils.Stack[*Category[T]]{},
		walkFlag:  walkLeaf,
	}

	qs.walk(qs.Contents, &ctx)
	return ctx.LeafCategories
}

func (qs *QuestionsDySchema[T]) WalkQuestions(walkFunc func(paths []*Category[T], question *T, appendable bool)) {
	ctx := &walkCtx[T]{
		pathStack: &utils.Stack[*Category[T]]{},
		walkFunc:  walkFunc,
	}
	qs.walk(qs.Contents, ctx)
}

func (qs *QuestionsDySchema[T]) walk(contents []*Category[T], ctx *walkCtx[T]) {
	for _, cate := range contents {
		// 入栈
		ctx.pathStack.Push(cate)

		if ctx.pathStack.Len() == qs.Leaf {
			// build flat question
			codePath := ctx.getPaths()
			fullPath := strings.Join(codePath, ".")

			if ctx.withLeaf() {
				ctx.LeafCategories = append(ctx.LeafCategories, &LeafInfo{
					Path:          fullPath,
					CategoryBasic: &cate.CategoryBasic,
					Replace:       cate.Replace,
				})
			}

			withQuest := ctx.withQuestion()
			if withQuest || ctx.walkFunc != nil {
				for _, q := range cate.Questions {
					if withQuest {
						fq := &FlatQuestion[T]{
							FullPathString: strings.Join(codePath, "."),
							Path:           codePath,
							PathCate:       ctx.pathStack.GetContent(),
							Question:       q,
							Appendable:     ctx.isAppendable(),
						}
						ctx.flatQuestions = append(ctx.flatQuestions, fq)
					}
					if ctx.walkFunc != nil {
						ctx.walkFunc(ctx.pathStack.GetContent(), q, ctx.isAppendable())
					}
				}
			}

		} else {
			qs.walk(cate.Children, ctx)
		}
		// 出栈
		ctx.pathStack.Pop()
	}
}

// walkCtx 的含义和作用是什么
type walkCtx[T any] struct {
	pathStack      *utils.Stack[*Category[T]]
	flatQuestions  []*FlatQuestion[T] // 问题列表
	LeafCategories []*LeafInfo        // walkLeaf类型 需要填充这个参数
	walkFlag       int
	walkFunc       func(paths []*Category[T], question *T, appendable bool)
}

type FlatQuestion[T any] struct {
	FullPathString string         `json:"fullPathString"` // 全路径
	Path           []string       `json:"path"`           // 路径上的每一个节点的id
	PathCate       []*Category[T] `json:"pathCate"`
	Question       *T             `json:"question"`   // 问题本身
	Appendable     bool           `json:"appendable"` // 是否可追加
}

// getPaths 获取路径上的每一个节点的id（code）
func (w *walkCtx[T]) getPaths() []string {
	catePath := w.pathStack.GetContent()
	var paths []string
	for _, cate := range catePath {
		paths = append(paths, cate.Code)
	}
	return paths
}

func (w *walkCtx[T]) withLeaf() bool {
	return w.walkFlag&walkLeaf > 0
}

func (w *walkCtx[T]) withQuestion() bool {
	return w.walkFlag&walkQuestion > 0
}
func (w *walkCtx[T]) isAppendable() bool {
	catePath := w.pathStack.GetContent()
	for _, cate := range catePath {
		if cate.Appendable {
			return true
		}
	}
	return false
}
