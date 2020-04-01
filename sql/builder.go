package sql

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

type Builder struct {
	SelectColumns []string
	WhereParams   []WhereParams
	OrderParams   []OrderParams
	Pages         *Pages
	Db            *gorm.DB
}

type WhereParams struct {
	Query interface{}
	Args  []interface{}
}

type OrderParams struct {
	Column string
	Asc    bool
}

func NewBuilder(selectColumns ...string) *Builder {
	s := &Builder{}
	if len(selectColumns) > 0 {
		s.SelectColumns = append(s.SelectColumns, selectColumns...)
	}
	return s
}

func (b *Builder) DB(db *gorm.DB) *Builder {
	b.Db = db
	return b
}

// where(" = ?", args)
// where(" <> ?", args)
// where(" > ?", args)
// where(" >= ?", args)
// where(" < ?", args)
// where(" <= ?", args)
// where(" in (?)", args)
// where(" LIKE ?", args)
func (b *Builder) Where(query string, args ...interface{}) *Builder {
	b.WhereParams = append(b.WhereParams, WhereParams{
		Query: query,
		Args:  args,
	})
	return b
}

func (b *Builder) Asc(column string) *Builder {
	b.OrderParams = append(b.OrderParams, OrderParams{
		Column: column,
		Asc:    true,
	})
	return b
}

func (b *Builder) Desc(column string) *Builder {
	b.OrderParams = append(b.OrderParams, OrderParams{
		Column: column,
		Asc:    false,
	})
	return b
}

func (b *Builder) Limit(limit int) *Builder {
	b.Page(1, limit)
	return b
}

func (b *Builder) Page(page, perPage int) *Builder {
	if b.Pages == nil {
		b.Pages = &Pages{
			Page:    page,
			PerPage: perPage,
		}
	} else {
		b.Pages.Page = page
		b.Pages.PerPage = perPage
	}
	return b
}

func (b *Builder) Build(count bool) *gorm.DB {

	if count && len(b.SelectColumns) > 0 {
		b.Db = b.Db.Select(b.SelectColumns)
	}

	if len(b.WhereParams) > 0 {
		for _, param := range b.WhereParams {
			b.Db = b.Db.Where(param.Query, param.Args...)
		}
	}

	if !count && len(b.OrderParams) > 0 {
		for _, order := range b.OrderParams {
			if order.Asc {
				b.Db = b.Db.Order(order.Column + " ASC")
			} else {
				b.Db = b.Db.Order(order.Column + " DESC")
			}
		}
	}

	if !count && b.Pages != nil && b.Pages.PerPage > 0 {
		b.Db = b.Db.Limit(b.Pages.PerPage)
	}

	if !count && b.Pages != nil && b.Pages.Offset() > 0 {
		b.Db = b.Db.Offset(b.Pages.Offset())
	}

	return b.Db
}

func SqlNullString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  len(value) > 0,
	}
}
