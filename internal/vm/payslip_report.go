package vm

import (
	"payslip-generator-service/internal/entity"
	ulid "payslip-generator-service/pkg/database/gorm"
)

// PayslipReportEmployee represents employee data in a payslip report
// swagger:model PayslipReportEmployee
type PayslipReportEmployee struct {
	// Unique identifier of the employee
	// example: "01HXYZ123456789ABCDEFGHIJK"
	EmployeeID ulid.ULID `json:"id"`

	// Username of the employee
	// example: "john.doe"
	EmployeeUsername string `json:"username"`

	// Employee's base salary amount
	// example: 5000000
	BasicSalary int `json:"basic_salary"`

	// Calculated salary for the period
	// example: 4500000
	Salary int `json:"salary"`

	// Final take-home pay for the employee
	// example: 4700000
	TakeHomePay int `json:"take_home_pay"`
}

// PayslipReport represents a comprehensive report of all employee payslips
// swagger:model PayslipReport
type PayslipReport struct {
	// List of employee payslip data
	Employees []PayslipReportEmployee `json:"employees"`

	// Total basic salary for all employees
	// example: 15000000
	TotalBasicSalary int `json:"total_basic_salary"`

	// Total calculated salary for all employees
	// example: 13500000
	TotalSalary int `json:"total_salary"`

	// Total take-home pay for all employees
	// example: 14100000
	TotalTakeHomePay int `json:"total_take_home_pay"`
}

// CreatePayslipReportProps represents the properties needed to create a new payslip report
// swagger:model CreatePayslipReportProps
type CreatePayslipReportProps struct {
	// List of all employees
	Employees []entity.Employee
	// List of all payslips
	Payslips []Payslip
}

func NewPayslipReport(props *CreatePayslipReportProps) *PayslipReport {

	employees := make([]PayslipReportEmployee, 0)
	totalBasicSalary := 0
	totalSalary := 0
	totalTakeHomePay := 0
	for _, employee := range props.Employees {
		var payslip Payslip
		for _, p := range props.Payslips {
			if p.EmployeeID == employee.ID {
				payslip = p
				break
			}
		}

		employees = append(employees, PayslipReportEmployee{
			EmployeeID:       employee.ID,
			EmployeeUsername: employee.Username,
			BasicSalary:      payslip.BasicSalary,
			Salary:           payslip.Salary,
			TakeHomePay:      payslip.TakeHomePay,
		})

		totalBasicSalary += payslip.BasicSalary
		totalSalary += payslip.Salary
		totalTakeHomePay += payslip.TakeHomePay
	}

	return &PayslipReport{
		Employees:        employees,
		TotalBasicSalary: totalBasicSalary,
		TotalSalary:      totalSalary,
		TotalTakeHomePay: totalTakeHomePay,
	}
}
