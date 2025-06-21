package vm

import (
	"payslip-generator-service/internal/entity"
	ulid "payslip-generator-service/pkg/database/gorm"
)

// reimbursementProps represents reimbursement summary data
// swagger:model reimbursementProps
type reimbursementProps struct {
	// Total number of reimbursement items
	// example: 3
	TotalItem int `json:"total_item"`

	// Total amount of all reimbursements
	// example: 450000
	TotalAmount int `json:"total_amount"`

	// List of reimbursement records
	Reimbursements []entity.Reimbursement `json:"reimbursements"`
}

// overtimeProps represents overtime summary data
// swagger:model overtimeProps
type overtimeProps struct {
	// Total number of overtime records
	// example: 5
	TotalItem int `json:"total_item"`

	// Total overtime pay amount
	// example: 250000
	TotalAmount int `json:"total_amount"`

	// Total overtime hours worked
	// example: 10
	TotalHours int `json:"total_hours"`

	// List of overtime records
	Overtimes []entity.Overtime `json:"overtimes"`
}

// Payslip represents a comprehensive payslip for an employee
// swagger:model Payslip
type Payslip struct {
	// Unique identifier of the employee
	// example: "01HXYZ123456789ABCDEFGHIJK"
	EmployeeID ulid.ULID `json:"employee_id"`

	// List of attendance records for the period
	Attendances []entity.Attendance `json:"attendances"`

	// Overtime summary and details
	Overtime overtimeProps `json:"overtime"`

	// Reimbursement summary and details
	Reimbursement reimbursementProps `json:"reimbursement"`

	// Employee's base salary amount
	// example: 5000000
	BasicSalary int `json:"basic_salary"`

	// Calculated salary for the period based on attendance
	// example: 4500000
	Salary int `json:"salary"`

	// Final take-home pay after deductions and additions
	// example: 4700000
	TakeHomePay int `json:"take_home_pay"`
}

// CreatePayslipProps represents the properties needed to create a new payslip
// swagger:model CreatePayslipProps
type CreatePayslipProps struct {
	// Unique identifier of the employee
	EmployeeID ulid.ULID
	// List of attendance records
	Attendance []entity.Attendance
	// List of overtime records
	Overtime []entity.Overtime
	// List of reimbursement records
	Reimbursement []entity.Reimbursement
	// Payroll period information
	PayrollPeriod entity.PayrollPeriod
	// Employee's base salary
	Salary int
}

func NewPayslip(props *CreatePayslipProps) *Payslip {
	maxSubmittedAt := props.PayrollPeriod.ProcessedAt

	// filter attendance (created_at <= maxSubmitedAt)
	attendances := make([]entity.Attendance, 0)
	for _, a := range props.Attendance {
		if a.CreatedAt.Before(*maxSubmittedAt) {
			attendances = append(attendances, a)
		}
	}

	totalDaysInPeriod := props.PayrollPeriod.GetDurationInDays()
	totalAttendance := min(len(attendances), totalDaysInPeriod) // get the minimum between the total attendance and the total days in period
	salaryPerDay := props.Salary / totalDaysInPeriod
	salaryPerHour := salaryPerDay / 8 // 8 hours per day
	salaryInPeriod := salaryPerDay * totalAttendance

	// filter overtime (created_at <= maxSubmitedAt)
	overtimes := make([]entity.Overtime, 0)
	totalAmountOvertime := 0
	totalHoursOvertime := 0
	for _, o := range props.Overtime {
		if o.CreatedAt.Before(*maxSubmittedAt) {
			overtimes = append(overtimes, o)
			totalAmountOvertime += o.TotalHours * (salaryPerHour * 2) // 2x salary per hour
			totalHoursOvertime += o.TotalHours
		}
	}

	// filter reimbursement (created_at <= maxSubmitedAt)
	reimbursements := make([]entity.Reimbursement, 0)
	totalAmountReimbursement := 0
	for _, r := range props.Reimbursement {
		if r.CreatedAt.Before(*maxSubmittedAt) {
			reimbursements = append(reimbursements, r)
			totalAmountReimbursement += r.Amount
		}
	}

	return &Payslip{
		EmployeeID:  props.EmployeeID,
		Attendances: attendances,
		Overtime: overtimeProps{
			TotalItem:   len(overtimes),
			TotalAmount: totalAmountOvertime,
			TotalHours:  totalHoursOvertime,
			Overtimes:   overtimes,
		},
		Reimbursement: reimbursementProps{
			TotalItem:      len(reimbursements),
			TotalAmount:    totalAmountReimbursement,
			Reimbursements: reimbursements,
		},
		BasicSalary: props.Salary,
		Salary:      salaryInPeriod,
		TakeHomePay: salaryInPeriod - totalAmountReimbursement + totalAmountOvertime,
	}
}
