package model

import "gorm.io/gorm"

type WebResponse[T any] struct {
	Ok     bool          `json:"ok"`
	Data   T             `json:"data,omitempty"`
	Errors interface{}   `json:"errors,omitempty"`
	Paging *PageMetadata `json:"paging,omitempty"`
}

type WebResponseWithData[T any] struct {
	Ok     bool          `json:"ok"`
	Data   T             `json:"data"`
	Errors interface{}   `json:"errors,omitempty"`
	Paging *PageMetadata `json:"paging,omitempty"`
}

type PageResponse[T any] struct {
	Data         []T          `json:"data,omitempty"`
	PageMetadata PageMetadata `json:"paging,omitempty"`
}

type OrderBy struct {
	Column    string
	Direction OrderDirection
}

type OrderDirection string

const (
	OrderDirectionAsc  OrderDirection = "ASC"
	OrderDirectionDesc OrderDirection = "DESC"
)

type PaginationOptions struct {
	Page     int
	PageSize int
	Filter   *func(tx *gorm.DB) *gorm.DB
	Order    []OrderBy
}

type PageMetadata struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}
