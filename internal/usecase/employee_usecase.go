package usecase

import (
	"context"
	"errors"
	"fmt"
	"payslip-generator-service/internal/entity"
	"payslip-generator-service/internal/repository"
	ulid "payslip-generator-service/pkg/database/gorm"
	"payslip-generator-service/pkg/logger"

	"gorm.io/gorm"
)

type EmployeeUseCase struct {
	DB                 *gorm.DB
	Log                *logger.ContextLogger
	EmployeeRepository *repository.EmployeeRepository
}

func NewEmployeeUseCase(
	db *gorm.DB,
	log *logger.ContextLogger,
	employeeRepository *repository.EmployeeRepository,
) *EmployeeUseCase {
	return &EmployeeUseCase{
		DB:                 db,
		Log:                log,
		EmployeeRepository: employeeRepository,
	}
}

func (a *EmployeeUseCase) GetById(ctx context.Context, id ulid.ULID) (*entity.Employee, error) {
	method := "EmployeeUseCase.GetById"
	a.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")
	a.Log.WithContext(ctx).WithField("method", method).WithField("request", id).Debug("request")

	db := a.DB.WithContext(ctx)

	employee := new(entity.Employee)
	if err := a.EmployeeRepository.FindById(db, employee, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("employee/not-found")
		}
		panic(err)
	}

	a.Log.WithContext(ctx).WithField("method", method).Trace("[END]")
	return employee, nil
}

func (a *EmployeeUseCase) List(ctx context.Context) ([]entity.Employee, error) {
	method := "EmployeeUseCase.List"
	a.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")
	a.Log.WithContext(ctx).WithField("method", method).Debug("request")

	db := a.DB.WithContext(ctx)

	employees := make([]entity.Employee, 0)
	if err := a.EmployeeRepository.FindAll(db, &employees); err != nil {
		panic(err)
	}

	a.Log.WithContext(ctx).WithField("method", method).Trace("[END]")
	return employees, nil
}
