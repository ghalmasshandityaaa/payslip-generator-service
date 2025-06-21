package usecase

import (
	"context"
	"errors"
	"fmt"
	"payslip-generator-service/internal/entity"
	"payslip-generator-service/internal/repository"
	ulid "payslip-generator-service/pkg/database/gorm"

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
	method := "EmployeeUseCase.GetById"
	a.Log.Trace("[BEGIN] - ", method)
	a.Log.Debug("request - ", method, map[string]interface{}{"employeeID": employeeID})

	db := a.DB.WithContext(ctx)

	employee := new(entity.Employee)
	if err := a.EmployeeRepository.FindById(db, employee, employeeID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("employee/not-found")
		}
		panic(err)
	}

	a.Log.Trace("[END] - ", method)
	return employee, nil
}

func (a *EmployeeUseCase) List(ctx context.Context) ([]entity.Employee, error) {
	method := "EmployeeUseCase.List"
	a.Log.Trace("[BEGIN] - ", method)
	a.Log.Debug("request - ", method)

	db := a.DB.WithContext(ctx)

	employees := make([]entity.Employee, 0)
	if err := a.EmployeeRepository.FindAll(db, &employees); err != nil {
		panic(err)
	}

	a.Log.Trace("[END] - ", method)
	return employees, nil
}
