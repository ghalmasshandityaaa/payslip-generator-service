package entity

import (
	"time"

	"payslip-generator-service/pkg/database/gorm"

	"github.com/oklog/ulid/v2"
)

// Overtime represents an employee's overtime record
// swagger:model Overtime
type Overtime struct {
	// Unique identifier for the overtime record
	// example: "01HXYZ123456789ABCDEFGHIJK"
	ID gorm.ULID `json:"id" gorm:"column:id;type:ulid;primaryKey"`

	// Date of the overtime work
	// example: "2024-01-15T00:00:00Z"
	Date time.Time `json:"date" gorm:"column:date;type:date;not null"`

	// Total hours of overtime (1-3 hours maximum)
	// example: 2
	TotalHours int `json:"total_hours" gorm:"column:total_hours;type:integer;not null"`

	// Timestamp when the overtime record was created
	// example: "2024-01-15T08:00:00Z"
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`

	// ID of the employee who created the overtime record
	// example: "01HXYZ123456789ABCDEFGHIJK"
	CreatedBy gorm.ULID `json:"created_by" gorm:"column:created_by;type:ulid;not null"`

	// Timestamp when the overtime record was last updated
	// example: "2024-01-15T08:00:00Z"
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp with time zone"`

	// ID of the employee who last updated the overtime record
	// example: "01HXYZ123456789ABCDEFGHIJK"
	UpdatedBy *gorm.ULID `json:"updated_by" gorm:"column:updated_by;type:ulid"`

	// Relations
	// Employee who created the overtime record
	Creator *Employee `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	// Employee who last updated the overtime record
	Updater *Employee `json:"updater,omitempty" gorm:"foreignKey:UpdatedBy"`
}

// CreateOvertimeProps represents the properties needed to create a new overtime record
// swagger:model CreateOvertimeProps
type CreateOvertimeProps struct {
	// Date of the overtime work
	Date time.Time
	// Total hours of overtime
	TotalHours int
	// ID of the employee creating the overtime record
	CreatedBy gorm.ULID
}

func NewOvertime(props *CreateOvertimeProps) *Overtime {
	return &Overtime{
		ID:         gorm.ULID(ulid.Make()),
		Date:       props.Date,
		TotalHours: props.TotalHours,
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

// Update updates the overtime with new data
func (o *Overtime) Update(date time.Time, totalHours int, updatedBy gorm.ULID) {
	now := time.Now()
	o.Date = date
	o.TotalHours = totalHours
	o.UpdatedAt = &now
	updatedByULID := gorm.ULID(updatedBy)
	o.UpdatedBy = &updatedByULID
}
