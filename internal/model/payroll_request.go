package model

import (
	"payslip-generator-service/internal/entity"

	ulid "payslip-generator-service/pkg/database/gorm"
)

// ListPayrollPeriodRequest represents the request parameters for listing payroll periods
// swagger:model ListPayrollPeriodRequest
type ListPayrollPeriodRequest struct {
	// Page number for pagination (minimum 1)
	// required: false
	// minimum: 1
	// example: 1
	Page int `json:"page" validate:"min=1"`

	// Number of items per page (minimum 1)
	// required: false
	// minimum: 1
	// example: 10
	PageSize int `json:"size" validate:"min=1"`
}

// CreatePayrollPeriodRequest represents the request body for creating a new payroll period
// swagger:model CreatePayrollPeriodRequest
type CreatePayrollPeriodRequest struct {
	// Start date of the payroll period (YYYY-MM-DD format)
	// required: true
	// example: "2024-01-01"
	StartDate string `json:"start_date" validate:"required,is-valid-date"`

	// End date of the payroll period (YYYY-MM-DD format)
	// required: true
	// example: "2024-01-31"
	EndDate string `json:"end_date" validate:"required,is-valid-date"`
}

// ProcessPayrollRequest represents the request body for processing payroll
// swagger:model ProcessPayrollRequest
type ProcessPayrollRequest struct {
	// Unique identifier of the payroll period
	// required: true
	// example: "01HXYZ123456789ABCDEFGHIJK"
	PeriodID string `json:"period_id" validate:"required,ulid"`
}

// GetPayslipRequest represents the request parameters for retrieving payslip
// swagger:model GetPayslipRequest
type GetPayslipRequest struct {
	// Unique identifier of the payroll period
	// required: true
	// example: "01HXYZ123456789ABCDEFGHIJK"
	PeriodID string `json:"period_id" validate:"required,ulid"`
}

// GeneratePayslipRequest represents the request body for generating payslip
// swagger:model GeneratePayslipRequest
type GeneratePayslipRequest struct {
	// Unique identifier of the employee
	// required: true
	// example: "01HXYZ123456789ABCDEFGHIJK"
	EmployeeID ulid.ULID `json:"employee_id"`

	// Employee's base salary amount
	// required: true
	// minimum: 0
	// example: 5000000
	Salary int `json:"salary"`

	// Payroll period information
	// required: true
	Period entity.PayrollPeriod `json:"period"`
}
