package entity

import (
	"time"

	"payslip-generator-service/pkg/database/gorm"

	"github.com/oklog/ulid/v2"
)

// Attendance model
type Attendance struct {
	ID        gorm.ULID  `json:"id" gorm:"column:id;type:ulid;primaryKey"`
	StartTime time.Time  `json:"start_time" gorm:"column:start_time;type:timestamp with time zone;not null"`
	EndTime   time.Time  `json:"end_time" gorm:"column:end_time;type:timestamp with time zone;not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy gorm.ULID  `json:"created_by" gorm:"column:created_by;type:ulid;not null"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp with time zone"`
	UpdatedBy *gorm.ULID `json:"updated_by" gorm:"column:updated_by;type:ulid"`

	// Relations
	Creator *Employee `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Updater *Employee `json:"updater,omitempty" gorm:"foreignKey:UpdatedBy"`
}

type CreateAttendanceProps struct {
	StartTime time.Time
	EndTime   time.Time
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
