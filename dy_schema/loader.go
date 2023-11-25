package dy_schema

import "github.com/weslynmxdlovewyn/tree/utils"

type copyCtx[T any] struct {
	Leaf      int
	pathStack *utils.Stack[*Category[T]]
}

func TransferDySchema[T, V any](src *QuestionsDySchema[T],
	transFunc func(paths []*Category[T], v *T, destParent *Category[V], appendable bool) (*V, error)) (*QuestionsDySchema[V], error) {
	dest := &QuestionsDySchema[V]{
		Leaf:          src.Leaf,
		MaxQuestionId: src.MaxQuestionId,
	}

	cctx := &copyCtx[T]{
		Leaf:      src.Leaf,
		pathStack: &utils.Stack[*Category[T]]{},
	}

	for _, cate := range src.Contents {
		cp, err := copy(cctx, cate, nil, transFunc)
		if err != nil {
			return nil, err
		}
		dest.Contents = append(dest.Contents, cp)
	}
	return dest, nil
}

func copy[T, V any](cctx *copyCtx[T], src *Category[T], destParent *Category[V],
	transFunc func(paths []*Category[T], v *T, destParent *Category[V], appendable bool) (*V, error)) (*Category[V], error) {
	dest := &Category[V]{
		CategoryBasic: src.CategoryBasic,
		Appendable:    src.Appendable,
	}

	cctx.pathStack.Push(src)

	if cctx.Leaf == cctx.pathStack.Len() {
		paths := cctx.pathStack.GetContent()
		appendable := paths[cctx.Leaf-2].Appendable
		for _, p := range src.Questions {
			nobj, err := transFunc(paths, p, destParent, appendable)
			if err != nil {
				return nil, err
			}
			dest.Questions = append(dest.Questions, nobj)
		}
	} else {
		for _, cate := range src.Children {
			cp, err := copy(cctx, cate, dest, transFunc)
			if err != nil {
				return nil, err
			}
			dest.Children = append(dest.Children, cp)
		}
	}
	cctx.pathStack.Pop()

	return dest, nil
}
