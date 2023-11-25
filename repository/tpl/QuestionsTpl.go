package tpl

import (
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
)

var QuestionsTplFields = struct {
	Id         string
	Version    string
	Name       string
	Text       string
	Keywords   string
	PageLayout string
	IsSealed   string
	CreatedBy  string
	CreatedAt  string
	ModifiedBy string
	ModifiedAt string
}{
	"id",
	"version",
	"name",
	"text",
	"keywords",
	"page_layout",
	"is_sealed",
	"created_by",
	"created_at",
	"modified_by",
	"modified_at",
}

var QuestionsTplMeta = &daog.TableMeta[QuestionsTpl]{
	Table: "questions_tpl",
	Columns: []string{
		"id",
		"version",
		"name",
		"text",
		"keywords",
		"page_layout",
		"is_sealed",
		"created_by",
		"created_at",
		"modified_by",
		"modified_at",
	},
	AutoColumn: "id",
	LookupFieldFunc: func(columnName string, ins *QuestionsTpl, point bool) any {
		if "id" == columnName {
			if point {
				return &ins.Id
			}
			return ins.Id
		}
		if "version" == columnName {
			if point {
				return &ins.Version
			}
			return ins.Version
		}
		if "name" == columnName {
			if point {
				return &ins.Name
			}
			return ins.Name
		}
		if "text" == columnName {
			if point {
				return &ins.Text
			}
			return ins.Text
		}
		if "keywords" == columnName {
			if point {
				return &ins.Keywords
			}
			return ins.Keywords
		}
		if "page_layout" == columnName {
			if point {
				return &ins.PageLayout
			}
			return ins.PageLayout
		}
		if "is_sealed" == columnName {
			if point {
				return &ins.IsSealed
			}
			return ins.IsSealed
		}
		if "created_by" == columnName {
			if point {
				return &ins.CreatedBy
			}
			return ins.CreatedBy
		}
		if "created_at" == columnName {
			if point {
				return &ins.CreatedAt
			}
			return ins.CreatedAt
		}
		if "modified_by" == columnName {
			if point {
				return &ins.ModifiedBy
			}
			return ins.ModifiedBy
		}
		if "modified_at" == columnName {
			if point {
				return &ins.ModifiedAt
			}
			return ins.ModifiedAt
		}

		return nil
	},
}

var QuestionsTplDao daog.QuickDao[QuestionsTpl] = &struct {
	daog.QuickDao[QuestionsTpl]
}{
	daog.NewBaseQuickDao(QuestionsTplMeta),
}

type QuestionsTpl struct {
	Id         int64
	Version    string
	Name       string
	Text       string
	Keywords   ttypes.NilableString
	PageLayout ttypes.NilableString
	IsSealed   int8
	CreatedBy  int64
	CreatedAt  ttypes.NormalDatetime
	ModifiedBy int64
	ModifiedAt ttypes.NormalDatetime
}
