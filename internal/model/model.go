package model

type WebResponse[T any] struct {
	Ok     bool          `json:"ok"`
	Data   T             `json:"data,omitempty"`
	Errors interface{}   `json:"errors,omitempty"`
	Paging *PageMetadata `json:"paging,omitempty"`
}

type PageResponse[T any] struct {
	Data         []T          `json:"data,omitempty"`
	PageMetadata PageMetadata `json:"paging,omitempty"`
}

type PageMetadata struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}
