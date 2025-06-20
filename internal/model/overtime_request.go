package model

type CreateOvertimeRequest struct {
	Date       string `json:"date" validate:"required,is-valid-date"`
	TotalHours int    `json:"total_hours" validate:"required,min=1,max=3"`
}
