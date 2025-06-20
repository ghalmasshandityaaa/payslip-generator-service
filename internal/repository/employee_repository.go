package repository

import (
	"payslip-generator-service/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmployeeRepository struct {
	Repository[entity.Employee]
	Log *logrus.Logger
}

func NewEmployeeRepository(log *logrus.Logger) *EmployeeRepository {
	return &EmployeeRepository{
		Log: log,
	}
}

func (r *EmployeeRepository) GetByUsername(db *gorm.DB, employee *entity.Employee, username string) error {
	return db.Debug().Where("username = ?", username).Take(employee).Error
}
