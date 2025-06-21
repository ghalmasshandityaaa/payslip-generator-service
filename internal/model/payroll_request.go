package model

import (
	"payslip-generator-service/internal/entity"
	"time"

	ulid "payslip-generator-service/pkg/database/gorm"
)

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

type GeneratePayslipRequest struct {
	EmployeeID ulid.ULID            `json:"employee_id"`
	Salary     int                  `json:"salary"`
	StartDate  time.Time            `json:"start_date"`
	EndDate    time.Time            `json:"end_date"`
	Period     entity.PayrollPeriod `json:"period"`
}
