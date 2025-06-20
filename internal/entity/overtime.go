package entity

import (
	"time"

	"payslip-generator-service/pkg/database/gorm"

	"github.com/oklog/ulid/v2"
)

// OvertimeStatus represents the status of overtime
type OvertimeStatus string

const (
	OvertimeStatusPending  OvertimeStatus = "pending"
	OvertimeStatusApproved OvertimeStatus = "approved"
	OvertimeStatusRejected OvertimeStatus = "rejected"
)

// Overtime model
type Overtime struct {
	ID         gorm.ULID      `json:"id" gorm:"column:id;type:ulid;primaryKey"`
	Date       time.Time      `json:"date" gorm:"column:date;type:date;not null"`
	TotalHours int            `json:"total_hours" gorm:"column:total_hours;type:integer;not null"`
	Status     OvertimeStatus `json:"status" gorm:"column:status;type:overtime_status;not null;default:'pending'"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy  gorm.ULID      `json:"created_by" gorm:"column:created_by;type:ulid;not null"`
	UpdatedAt  *time.Time     `json:"updated_at" gorm:"column:updated_at;type:timestamp with time zone"`
	UpdatedBy  *gorm.ULID     `json:"updated_by" gorm:"column:updated_by;type:ulid"`

	// Relations
	Creator *Employee `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Updater *Employee `json:"updater,omitempty" gorm:"foreignKey:UpdatedBy"`
}

type CreateOvertimeProps struct {
	Date       time.Time
	TotalHours int
	CreatedBy  gorm.ULID
}

func NewOvertime(props *CreateOvertimeProps) *Overtime {
	return &Overtime{
		ID:         gorm.ULID(ulid.Make()),
		Date:       props.Date,
		TotalHours: props.TotalHours,
		Status:     OvertimeStatusPending,
		CreatedAt:  time.Now(),
		CreatedBy:  props.CreatedBy,
	}
}

func (o *Overtime) TableName() string {
	return "overtime"
}

// IsValidDuration checks if the overtime duration is within valid range (1-3 hours)
func (o *Overtime) IsValidDuration() bool {
	return o.TotalHours >= 1 && o.TotalHours <= 3
}

// IsPending checks if the overtime is still pending
func (o *Overtime) IsPending() bool {
	return o.Status == OvertimeStatusPending
}

// IsApproved checks if the overtime is approved
func (o *Overtime) IsApproved() bool {
	return o.Status == OvertimeStatusApproved
}

// IsRejected checks if the overtime is rejected
func (o *Overtime) IsRejected() bool {
	return o.Status == OvertimeStatusRejected
}

// IsToday checks if the overtime is for today
func (o *Overtime) IsToday() bool {
	return o.Date.Year() == time.Now().Year() &&
		o.Date.YearDay() == time.Now().YearDay()
}

// IsWeekday checks if the overtime is on a weekday
func (o *Overtime) IsWeekday() bool {
	weekday := o.Date.Weekday()
	// Sunday = 0, Saturday = 6
	return weekday != time.Sunday && weekday != time.Saturday
}

// Approve approves the overtime
func (o *Overtime) Approve(updatedBy gorm.ULID) {
	now := time.Now()
	o.Status = OvertimeStatusApproved
	o.UpdatedAt = &now
	updatedByULID := gorm.ULID(updatedBy)
	o.UpdatedBy = &updatedByULID
}

// Reject rejects the overtime
func (o *Overtime) Reject(updatedBy gorm.ULID) {
	now := time.Now()
	o.Status = OvertimeStatusRejected
	o.UpdatedAt = &now
	updatedByULID := gorm.ULID(updatedBy)
	o.UpdatedBy = &updatedByULID
}

// Update updates the overtime with new data
func (o *Overtime) Update(date time.Time, totalHours int, updatedBy gorm.ULID) {
	now := time.Now()
	o.Date = date
	o.TotalHours = totalHours
	o.UpdatedAt = &now
	updatedByULID := gorm.ULID(updatedBy)
	o.UpdatedBy = &updatedByULID
}
