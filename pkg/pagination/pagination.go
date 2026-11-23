// Package pagination provides the pagination and filtering helpers used by the application.
package pagination

import (
	"net/url"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type PagingParams[T any] struct {
	Host    string
	Sort    string
	Order   string
	Page    int
	Limit   int
	Filters T
}

type PagingLinks struct {
	First    string `json:"first"`
	Last     string `json:"last"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

type PagingInfo[T any] struct {
	Links      PagingLinks `json:"links"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
	Rows       T           `json:"rows"`
}

type CommonTimeFilters struct {
	DateField, CreatedAfter, UpdatedAfter, CreatedBefore, UpdatedBefore string
	CreatedAt, UpdatedAt, StartDate, EndDate                            string
}

type SortParams struct{ Sort, Order, DefaultOrderBy string }

func NewPaginationParams[T any](host, sort, order, page, limit, defaultSort string) PagingParams[T] {
	p, err := strconv.Atoi(page)
	if err != nil || p < 1 {
		p = 1
	}
	size, err := strconv.Atoi(limit)
	if err != nil || size < 1 {
		size = 10
	}
	if sort == "" {
		fields := strings.Fields(defaultSort)
		if len(fields) > 0 {
			sort = fields[0]
		}
		if len(fields) > 1 && order == "" {
			order = fields[1]
		}
	}
	return PagingParams[T]{Host: host, Sort: sort, Order: order, Page: p, Limit: size}
}

// NewOrderBy only accepts the todo columns supported for sorting.
func NewOrderBy(p SortParams) string {
	allowed := map[string]bool{"id": true, "title": true, "detail": true, "created_at": true, "updated_at": true}
	if !allowed[p.Sort] {
		return p.DefaultOrderBy
	}
	order := strings.ToUpper(p.Order)
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}
	return p.Sort + " " + order
}

func ApplyFilter(db *gorm.DB, column, value, operation string) *gorm.DB {
	if value == "" {
		return db
	}
	if operation == "contains" {
		return db.Where(column+" LIKE ?", "%"+value+"%")
	}
	return db.Where(column+" = ?", value)
}

func ApplyDatetimeFilters(db *gorm.DB, f CommonTimeFilters) *gorm.DB {
	for _, item := range []struct{ column, op, value string }{
		{"created_at", ">=", f.CreatedAfter}, {"updated_at", ">=", f.UpdatedAfter},
		{"created_at", "<=", f.CreatedBefore}, {"updated_at", "<=", f.UpdatedBefore},
		{"created_at", "=", f.CreatedAt}, {"updated_at", "=", f.UpdatedAt},
	} {
		if item.value != "" {
			db = db.Where(item.column+" "+item.op+" ?", item.value)
		}
	}
	field := f.DateField
	if field != "created_at" && field != "updated_at" {
		field = "created_at"
	}
	if f.StartDate != "" {
		db = db.Where(field+" >= ?", f.StartDate)
	}
	if f.EndDate != "" {
		db = db.Where(field+" <= ?", f.EndDate)
	}
	return db
}

func Paginate[F, T any](p PagingParams[F], db *gorm.DB) (PagingInfo[T], error) {
	var result PagingInfo[T]
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 10
	}
	if err := db.Count(&result.Total).Error; err != nil {
		return result, err
	}
	result.Page, result.PageSize = p.Page, p.Limit
	result.TotalPages = int(result.Total / int64(p.Limit))
	if result.Total%int64(p.Limit) != 0 {
		result.TotalPages++
	}
	if err := db.Limit(p.Limit).Offset((p.Page - 1) * p.Limit).Find(&result.Rows).Error; err != nil {
		return result, err
	}
	link := func(page int) string {
		u, err := url.Parse(p.Host)
		if err != nil {
			return ""
		}
		q := u.Query()
		q.Set("page", strconv.Itoa(page))
		q.Set("limit", strconv.Itoa(p.Limit))
		if p.Sort != "" {
			q.Set("sort", p.Sort)
		}
		if p.Order != "" {
			q.Set("order", p.Order)
		}
		u.RawQuery = q.Encode()
		return u.String()
	}
	result.Links.First = link(1)
	last := result.TotalPages
	if last < 1 {
		last = 1
	}
	result.Links.Last = link(last)
	if p.Page > 1 {
		result.Links.Previous = link(p.Page - 1)
	}
	if p.Page < result.TotalPages {
		result.Links.Next = link(p.Page + 1)
	}
	return result, nil
}
