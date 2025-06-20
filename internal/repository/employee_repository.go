package repository

import (
	"payslip-generator-service/internal/entity"

	"github.com/sirupsen/logrus"
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
