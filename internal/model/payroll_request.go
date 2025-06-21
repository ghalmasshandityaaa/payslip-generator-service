package model

type ListPayrollPeriodRequest struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"size" validate:"min=1"`
}

type CreatePayrollPeriodRequest struct {
	StartDate string `json:"start_date" validate:"required,is-valid-date"`
	EndDate   string `json:"end_date" validate:"required,is-valid-date"`
}

type ProcessPayrollRequest struct {
	PeriodID string `json:"period_id" validate:"required,ulid"`
}

type GetPayslipRequest struct {
	PeriodID string `json:"period_id" validate:"required,ulid"`
}
