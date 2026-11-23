package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go-hexagonal-example/pkg/pagination"
)

type PaginationResponse struct {
	Links      pagination.PagingLinks `json:"links"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}

type APIPaginationResponse[T interface{}] struct {
	StatusCode    string             `json:"status_code"`
	StatusMessage string             `json:"status_message"`
	Data          T                  `json:"data"`
	Pagination    PaginationResponse `json:"pagination"`
}

func GetAPIEndpoint(c *fiber.Ctx) string {
	return fmt.Sprintf("%s://%s%s", c.Protocol(), c.Hostname(), c.OriginalURL())
}

func InitPagination[FilterType interface{}](c *fiber.Ctx) pagination.PagingParams[FilterType] {
	host := GetAPIEndpoint(c) // Assuming utils.GetAPIEndpoint is a function you have
	defaultSort := "created_at desc"
	sort := c.Query("sort")
	order := c.Query("order")
	page := c.Query("page")
	limit := c.Query("limit")
	return pagination.NewPaginationParams[FilterType](host, sort, order, page, limit, defaultSort)
}

func GetPaginationInfo[T interface{}](payload pagination.PagingInfo[T]) PaginationResponse {
	return PaginationResponse{
		Links:      payload.Links,
		Total:      payload.Total,
		Page:       payload.Page,
		PageSize:   payload.PageSize,
		TotalPages: payload.TotalPages,
	}
}

func NewPaginationResponse[T interface{}](c *fiber.Ctx, message string, data pagination.PagingInfo[T]) error {
	response := APIPaginationResponse[T]{
		StatusCode:    "200",
		StatusMessage: message,
		Data:          data.Rows,
		Pagination:    GetPaginationInfo[T](data),
	}
	return c.Status(200).JSON(response)
}
