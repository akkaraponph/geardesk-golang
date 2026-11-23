package core

import (
	"errors"
	"go-hexagonal-example/examples/filters"

	"go-hexagonal-example/examples/pkg/pagination"
)

type (
	ISrvInfs interface {
		GetTodoes(p pagination.PagingParams[filters.FTodo]) (pagination.PagingInfo[[]TodoDto], error)
	}
	srvDeps struct {
		todoRepo IRepoInfs
	}
)

func NewServices(
	todoRepo IRepoInfs,
) ISrvInfs {
	return &srvDeps{
		todoRepo: todoRepo,
	}
}

// GetTodoes implements ISrvInfs.
func (s *srvDeps) GetTodoes(p pagination.PagingParams[filters.FTodo]) (pagination.PagingInfo[[]TodoDto], error) {
	pgR, err := s.todoRepo.GetTodoes(p)
	if err != nil {
		return pagination.PagingInfo[[]TodoDto]{}, err
	}
	if pgR == nil {
		err := errors.New("no data found")
		return pagination.PagingInfo[[]TodoDto]{}, err
	}

	var resp []TodoDto
	for _, v := range pgR.Rows {
		resp = append(resp, TodoDto{
			ID:        v.ID,
			Title:     v.Title,
			Detail:    v.Detail,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return pagination.PagingInfo[[]TodoDto]{
		Links:      pgR.Links,
		Total:      pgR.Total,
		Page:       pgR.Page,
		PageSize:   pgR.PageSize,
		TotalPages: pgR.TotalPages,
		Rows:       resp,
	}, nil
}
