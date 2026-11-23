package core

import (
	"go-hexagonal-example/examples/filters"
	"go-hexagonal-example/examples/models"

	"go-hexagonal-example/examples/pkg/pagination"
	"gorm.io/gorm"
)

type (
	IRepoInfs interface {
		GetTodoes(p pagination.PagingParams[filters.FTodo]) (*pagination.PagingInfo[[]models.TodoModel], error)
	}
	repoDeps struct {
		db *gorm.DB
	}
)

func NewRepository(db *gorm.DB) IRepoInfs {
	return &repoDeps{db: db}
}

// GetTodoes implements IRepoInfs.
func (r *repoDeps) GetTodoes(p pagination.PagingParams[filters.FTodo]) (*pagination.PagingInfo[[]models.TodoModel], error) {

	fp := p.Filters
	q := r.db.Model(&models.TodoModel{})
	q = pagination.ApplyFilter(q, "title", fp.Title, "contains")

	q = pagination.ApplyDatetimeFilters(q, pagination.CommonTimeFilters{
		DateField:     fp.CommonTimeFilters.DateField,
		CreatedAfter:  fp.CommonTimeFilters.CreatedAfter,
		UpdatedAfter:  fp.CommonTimeFilters.UpdatedAfter,
		CreatedBefore: fp.CommonTimeFilters.CreatedBefore,
		UpdatedBefore: fp.CommonTimeFilters.UpdatedBefore,
		CreatedAt:     fp.CommonTimeFilters.CreatedAt,
		UpdatedAt:     fp.CommonTimeFilters.UpdatedAt,
		StartDate:     fp.CommonTimeFilters.StartDate,
		EndDate:       fp.CommonTimeFilters.EndDate,
	})
	orderBy := pagination.NewOrderBy(pagination.SortParams{
		Sort:           p.Sort,
		Order:          p.Order,
		DefaultOrderBy: "updated_at DESC",
	})
	q = q.Order(orderBy)
	pgR, err := pagination.Paginate[filters.FTodo, []models.TodoModel](p, q)
	if err != nil {
		return nil, err
	}
	return &pgR, nil
}
