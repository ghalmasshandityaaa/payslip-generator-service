package model

type ListPayrollPeriodRequest struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"size" validate:"min=1"`
}

type CreatePayrollPeriodResponse struct{}
