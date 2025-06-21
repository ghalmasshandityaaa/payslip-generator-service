package entity

import (
	"time"

	"payslip-generator-service/pkg/database/gorm"

	"github.com/oklog/ulid/v2"
)

// PayrollPeriod represents a payroll period in the system
// swagger:model PayrollPeriod
type PayrollPeriod struct {
	// Unique identifier for the payroll period
	// example: "01HXYZ123456789ABCDEFGHIJK"
	ID gorm.ULID `json:"id" gorm:"column:id;type:ulid;primaryKey"`

	// Start date of the payroll period
	// example: "2024-01-01T00:00:00Z"
	StartDate time.Time `json:"start_date" gorm:"column:start_date;type:date;not null"`

	// End date of the payroll period
	// example: "2024-01-31T00:00:00Z"
	EndDate time.Time `json:"end_date" gorm:"column:end_date;type:date;not null"`

	// Timestamp when the payroll was processed
	// example: "2024-01-31T23:59:59Z"
	ProcessedAt *time.Time `json:"processed_at" gorm:"column:processed_at;type:timestamp with time zone"`

	// ID of the employee who processed the payroll
	// example: "01HXYZ123456789ABCDEFGHIJK"
	ProcessedBy *gorm.ULID `json:"processed_by" gorm:"column:processed_by;type:ulid"`

	// Timestamp when the payroll period was created
	// example: "2024-01-01T08:00:00Z"
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`

	// ID of the employee who created the payroll period
	// example: "01HXYZ123456789ABCDEFGHIJK"
	CreatedBy gorm.ULID `json:"created_by" gorm:"column:created_by;type:ulid;not null"`

	// Timestamp when the payroll period was last updated
	// example: "2024-01-15T08:00:00Z"
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp with time zone"`

	// ID of the employee who last updated the payroll period
	// example: "01HXYZ123456789ABCDEFGHIJK"
	UpdatedBy *gorm.ULID `json:"updated_by" gorm:"column:updated_by;type:ulid"`

	// Relations
	// Employee who created the payroll period
	Creator *Employee `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	// Employee who last updated the payroll period
	Updater *Employee `json:"updater,omitempty" gorm:"foreignKey:UpdatedBy"`
}

// CreatePayrollPeriodProps represents the properties needed to create a new payroll period
// swagger:model CreatePayrollPeriodProps
type CreatePayrollPeriodProps struct {
	// Start date of the payroll period
	StartDate time.Time
	// End date of the payroll period
	EndDate time.Time
	// ID of the employee creating the payroll period
	CreatedBy gorm.ULID
}

func NewPayrollPeriod(props *CreatePayrollPeriodProps) *PayrollPeriod {
	return &PayrollPeriod{
		ID:        gorm.ULID(ulid.Make()),
		StartDate: props.StartDate,
		EndDate:   props.EndDate,
		CreatedAt: time.Now(),
		CreatedBy: props.CreatedBy,
	}
}

func (p *PayrollPeriod) TableName() string {
	return "payroll_period"
}

// GetDuration returns the duration of the payroll period
func (p *PayrollPeriod) GetDuration() time.Duration {
	endOfDay := time.Date(p.EndDate.Year(), p.EndDate.Month(), p.EndDate.Day(), 23, 59, 59, 999999999, p.EndDate.Location())
	return endOfDay.Sub(p.StartDate)
}

// GetDurationInDays returns the duration in days
func (p *PayrollPeriod) GetDurationInDays() int {
	return int((p.GetDuration() + time.Millisecond).Hours() / 24)
}

// IsValidDateRange checks if the start date is before the end date
func (p *PayrollPeriod) IsValidDateRange() bool {
	return p.StartDate.Before(p.EndDate) || p.StartDate.Equal(p.EndDate)
}

// Process marks the payroll as processed
func (p *PayrollPeriod) Process(processedBy gorm.ULID) {
	now := time.Now()
	p.ProcessedAt = &now
	p.ProcessedBy = &processedBy
}

// Update updates the payroll with new data
func (p *PayrollPeriod) Update(startDate, endDate time.Time, updatedBy gorm.ULID) {
	now := time.Now()
	p.StartDate = startDate
	p.EndDate = endDate
	p.UpdatedAt = &now
	p.UpdatedBy = &updatedBy
}

// IsCurrentPeriod checks if the payroll period includes the current date
func (p *PayrollPeriod) IsCurrentPeriod() bool {
	now := time.Now()
	return (p.StartDate.Equal(now) || p.StartDate.Before(now)) &&
		(p.EndDate.Equal(now) || p.EndDate.After(now))
}

// IsFuturePeriod checks if the payroll period is in the future
func (p *PayrollPeriod) IsFuturePeriod() bool {
	now := time.Now()
	return p.StartDate.After(now)
}

// IsPastPeriod checks if the payroll period is in the past
func (p *PayrollPeriod) IsPastPeriod() bool {
	now := time.Now()
	return p.EndDate.Before(now)
}

// IsProcessed checks if the payroll period is processed
func (p *PayrollPeriod) IsProcessed() bool {
	return p.ProcessedAt != nil && p.ProcessedBy != nil
}
