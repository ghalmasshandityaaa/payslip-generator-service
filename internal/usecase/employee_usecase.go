package usecase

import (
	"context"
	"payslip-generator-service/internal/entity"
	"payslip-generator-service/internal/repository"

	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmployeeUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	EmployeeRepository *repository.EmployeeRepository
}

func NewEmployeeUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	employeeRepository *repository.EmployeeRepository,
) *EmployeeUseCase {
	return &EmployeeUseCase{
		DB:                 db,
		Log:                logger,
		EmployeeRepository: employeeRepository,
	}
}

func (a *EmployeeUseCase) GetById(ctx context.Context, employeeID ulid.ULID) (*entity.Employee, error) {
	return nil, nil
}
