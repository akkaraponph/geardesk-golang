package filters

import "go-hexagonal-example/examples/pkg/pagination"

type FTodo struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	pagination.CommonTimeFilters
}
