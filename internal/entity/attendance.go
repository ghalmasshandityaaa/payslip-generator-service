package entity

import (
	"time"

	"payslip-generator-service/pkg/database/gorm"

	"github.com/oklog/ulid/v2"
)

// Attendance represents an employee's attendance record
// swagger:model Attendance
type Attendance struct {
	// Unique identifier for the attendance record
	// example: "01HXYZ123456789ABCDEFGHIJK"
	ID gorm.ULID `json:"id" gorm:"column:id;type:ulid;primaryKey"`

	// Start time of the work shift
	// example: "2024-01-15T08:00:00Z"
	StartTime time.Time `json:"start_time" gorm:"column:start_time;type:timestamp with time zone;not null"`

	// End time of the work shift
	// example: "2024-01-15T17:00:00Z"
	EndTime time.Time `json:"end_time" gorm:"column:end_time;type:timestamp with time zone;not null"`

	// Timestamp when the attendance record was created
	// example: "2024-01-15T08:00:00Z"
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`

	// ID of the employee who created the attendance record
	// example: "01HXYZ123456789ABCDEFGHIJK"
	CreatedBy gorm.ULID `json:"created_by" gorm:"column:created_by;type:ulid;not null"`

	// Timestamp when the attendance record was last updated
	// example: "2024-01-15T08:00:00Z"
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp with time zone"`

	// ID of the employee who last updated the attendance record
	// example: "01HXYZ123456789ABCDEFGHIJK"
	UpdatedBy *gorm.ULID `json:"updated_by" gorm:"column:updated_by;type:ulid"`

	// Relations
	// Employee who created the attendance record
	Creator *Employee `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	// Employee who last updated the attendance record
	Updater *Employee `json:"updater,omitempty" gorm:"foreignKey:UpdatedBy"`
}

// CreateAttendanceProps represents the properties needed to create a new attendance record
// swagger:model CreateAttendanceProps
type CreateAttendanceProps struct {
	// Start time of the work shift
	StartTime time.Time
	// End time of the work shift
	EndTime time.Time
	// ID of the employee creating the attendance record
	CreatedBy gorm.ULID
}

func NewAttendance(props *CreateAttendanceProps) *Attendance {
	return &Attendance{
		ID:        gorm.ULID(ulid.Make()),
		StartTime: props.StartTime,
		EndTime:   props.EndTime,
		CreatedAt: time.Now(),
		CreatedBy: props.CreatedBy,
	}
}

func (a *Attendance) TableName() string {
	return "attendance"
}

// GetDuration returns the duration of the attendance
func (a *Attendance) GetDuration() time.Duration {
	return a.EndTime.Sub(a.StartTime)
}

// GetDurationInHours returns the duration in hours
func (a *Attendance) GetDurationInHours() float64 {
	return a.GetDuration().Hours()
}

func (a *Attendance) IsToday() bool {
	return a.StartTime.Year() == time.Now().Year() &&
		a.StartTime.YearDay() == time.Now().YearDay() &&
		a.EndTime.Year() == time.Now().Year() &&
		a.EndTime.YearDay() == time.Now().YearDay()
}

// IsSameDay checks if the attendance spans the same day
func (a *Attendance) IsSameDay() bool {
	return a.StartTime.Year() == a.EndTime.Year() &&
		a.StartTime.YearDay() == a.EndTime.YearDay()
}

// IsWeekday checks if the attendance is on a weekday
func (a *Attendance) IsWeekday() bool {
	startWeekday := a.StartTime.Weekday()
	endWeekday := a.EndTime.Weekday()

	// Sunday = 0, Saturday = 6
	return startWeekday != time.Sunday && startWeekday != time.Saturday &&
		endWeekday != time.Sunday && endWeekday != time.Saturday
}

// endTime must be greater than startTime
func (a *Attendance) IsEndTimeGreaterThanStartTime() bool {
	return a.EndTime.After(a.StartTime)
}

// Update updates the attendance with new times
func (a *Attendance) Update(startTime, endTime time.Time, updatedBy gorm.ULID) {
	now := time.Now()
	a.StartTime = startTime
	a.EndTime = endTime
	a.UpdatedAt = &now
	updatedByULID := gorm.ULID(updatedBy)
	a.UpdatedBy = &updatedByULID
}
