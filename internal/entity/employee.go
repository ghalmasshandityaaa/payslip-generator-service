package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

// Employee model
type Employee struct {
	ID        ulid.ULID  `json:"id" gorm:"column:id;type:ulid;primaryKey"`
	Username  string     `json:"username" gorm:"column:username;size:50;not null;unique"`
	Password  string     `json:"-" gorm:"column:password;type:text;not null"`
	Salary    int        `json:"salary" gorm:"column:salary;type:integer;not null"`
	IsAdmin   bool       `json:"is_admin" gorm:"column:is_admin;type:boolean;not null;default:false"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp with time zone"`
}

type CreateEmployeeProps struct {
	Username string
	Password string
	Salary   int
	IsAdmin  bool
}

func NewEmployee(props *CreateEmployeeProps) *Employee {
	return &Employee{
		ID:        ulid.Make(),
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
