package entity

import (
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

type Gender string

const (
	MALE   Gender = "MALE"
	FEMALE Gender = "FEMALE"
)

// Employee model
type Employee struct {
	ID          ulid.ULID `json:"id" gorm:"column:id;type:ulid;primaryKey"`
	FirstName   string    `json:"first_name" gorm:"column:first_name;size:100;not null"`
	LastName    string    `json:"last_name" gorm:"column:last_name;size:100;not null"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone_number;size:16;not null;unique"`
	BirthDate   time.Time `json:"birth_date" gorm:"column:birth_date;type:date;not null"`
	IsAdmin     bool      `json:"is_admin" gorm:"column:is_admin;type:boolean;not null;default:false"`
	CreatedAt   int64     `json:"-" gorm:"column:created_at;autoCreateTime:milli;not null"`
	UpdatedAt   int64     `json:"-" gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli;not null"`
}

type CreateEmployeeProps struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	BirthDate   time.Time
	Gender      Gender
	IsAdmin     bool
}

func NewEmployee(user *CreateEmployeeProps) *Employee {
	return &Employee{
		ID:          ulid.Make(),
		FirstName:   strings.ToLower(user.FirstName),
		LastName:    strings.ToLower(user.LastName),
		PhoneNumber: user.PhoneNumber,
		BirthDate:   user.BirthDate,
		IsAdmin:     user.IsAdmin,
		CreatedAt:   time.Now().UnixMilli(),
		UpdatedAt:   time.Now().UnixMilli(),
	}
}

func (u *Employee) TableName() string {
	return "user"
}
