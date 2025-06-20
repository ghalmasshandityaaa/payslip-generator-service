package model

import "payslip-generator-service/pkg/database/gorm"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
)

type Auth struct {
	ID      gorm.ULID
	IsAdmin bool
}
