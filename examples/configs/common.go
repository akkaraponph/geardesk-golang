package configs

import "go-hexagonal-example/examples/pkg/pagination"

type PagingResponse[T interface{}] struct {
	Code          int                      `json:"code"`
	StatusCode    int                      `json:"status_code"`
	StatusMessage string                   `json:"status_message"`
	Data          interface{}              `json:"data"`
	Pagination    pagination.PagingInfo[T] `json:"pagination"`
}
