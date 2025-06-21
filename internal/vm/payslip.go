package vm

import (
	"payslip-generator-service/internal/entity"
	ulid "payslip-generator-service/pkg/database/gorm"
)

type reimbursementProps struct {
	TotalItem      int                    `json:"total_item"`
	TotalAmount    int                    `json:"total_amount"`
	Reimbursements []entity.Reimbursement `json:"reimbursements"`
}

type overtimeProps struct {
	TotalItem   int               `json:"total_item"`
	TotalAmount int               `json:"total_amount"`
	TotalHours  int               `json:"total_hours"`
	Overtimes   []entity.Overtime `json:"overtimes"`
}

// Payslip model
type Payslip struct {
	EmployeeID    ulid.ULID           `json:"employee_id"`
	Attendances   []entity.Attendance `json:"attendances"`
	Overtime      overtimeProps       `json:"overtime"`
	Reimbursement reimbursementProps  `json:"reimbursement"`
	BasicSalary   int                 `json:"basic_salary"`
	Salary        int                 `json:"salary"`
	TakeHomePay   int                 `json:"take_home_pay"`
}

type CreatePayslipProps struct {
	EmployeeID    ulid.ULID
	Attendance    []entity.Attendance
	Overtime      []entity.Overtime
	Reimbursement []entity.Reimbursement
	PayrollPeriod entity.PayrollPeriod
	Salary        int
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
