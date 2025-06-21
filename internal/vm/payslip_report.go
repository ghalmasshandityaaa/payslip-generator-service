package vm

import (
	"payslip-generator-service/internal/entity"
	ulid "payslip-generator-service/pkg/database/gorm"
)

type PayslipReportEmployee struct {
	EmployeeID       ulid.ULID `json:"id"`
	EmployeeUsername string    `json:"username"`
	BasicSalary      int       `json:"basic_salary"`
	Salary           int       `json:"salary"`
	TakeHomePay      int       `json:"take_home_pay"`
}

// PayslipReport model
type PayslipReport struct {
	Employees        []PayslipReportEmployee `json:"employees"`
	TotalBasicSalary int                     `json:"total_basic_salary"`
	TotalSalary      int                     `json:"total_salary"`
	TotalTakeHomePay int                     `json:"total_take_home_pay"`
}

type CreatePayslipReportProps struct {
	Employees []entity.Employee
	Payslips  []Payslip
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
