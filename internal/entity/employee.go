package entity

import (
	"time"

	"payslip-generator-service/pkg/database/gorm"

	"github.com/oklog/ulid/v2"
)

// Employee represents an employee in the system
// swagger:model Employee
type Employee struct {
	// Unique identifier for the employee
	// example: "01HXYZ123456789ABCDEFGHIJK"
	ID gorm.ULID `json:"id" gorm:"column:id;type:ulid;primaryKey"`

	// Username for authentication (unique)
	// example: "john.doe"
	Username string `json:"username" gorm:"column:username;size:50;not null;unique"`

	// Password hash (not exposed in JSON)
	Password string `json:"-" gorm:"column:password;type:text;not null"`

	// Employee's base salary amount
	// example: 5000000
	Salary int `json:"salary" gorm:"column:salary;type:integer;not null"`

	// Whether the employee has admin privileges
	// example: false
	IsAdmin bool `json:"is_admin" gorm:"column:is_admin;type:boolean;not null;default:false"`

	// Timestamp when the employee was created
	// example: "2024-01-15T08:00:00Z"
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`

	// Timestamp when the employee was last updated
	// example: "2024-01-15T08:00:00Z"
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp with time zone"`
}

// CreateEmployeeProps represents the properties needed to create a new employee
// swagger:model CreateEmployeeProps
type CreateEmployeeProps struct {
	// Username for the new employee
	Username string
	// Password for the new employee
	Password string
	// Base salary for the new employee
	Salary int
	// Whether the new employee should have admin privileges
	IsAdmin bool
}

func NewEmployee(props *CreateEmployeeProps) *Employee {
	return &Employee{
		ID:        gorm.ULID(ulid.Make()),
		Username:  props.Username,
		Password:  props.Password,
		Salary:    props.Salary,
		IsAdmin:   props.IsAdmin,
		CreatedAt: time.Now(),
	}
}

func (e *Employee) TableName() string {
	return "employee"
}
