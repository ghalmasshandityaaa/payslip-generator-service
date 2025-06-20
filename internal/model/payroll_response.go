package model

type CreatePayrollPeriodRequest struct {
	StartDate string `json:"start_date" validate:"required,is-valid-date"`
	EndDate   string `json:"end_date" validate:"required,is-valid-date"`
}
