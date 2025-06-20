package usecase

import (
	"context"
	"fmt"
	"payslip-generator-service/internal/entity"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/repository"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PayrollUseCase struct {
	DB            *gorm.DB
	Log           *logrus.Logger
	payrollPeriod *repository.PayrollPeriodRepository
}

func NewPayrollUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	payrollPeriodRepository *repository.PayrollPeriodRepository,
) *PayrollUseCase {
	return &PayrollUseCase{
		DB:            db,
		Log:           logger,
		payrollPeriod: payrollPeriodRepository,
	}
}

func (a *PayrollUseCase) ListPeriod(ctx context.Context, request *model.ListPayrollPeriodRequest) ([]entity.PayrollPeriod, int64, error) {
	method := "PayrollUseCase.ListPeriod"
	a.Log.Trace("[BEGIN] - ", method)
	a.Log.Debug("request - ", method, request)

	db := a.DB.WithContext(ctx)
	data, total, err := a.payrollPeriod.FindAllWithPagination(db, &model.PaginationOptions{
		Page:     request.Page,
		PageSize: request.PageSize,
		Order: []model.OrderBy{
			{
				Column:    "start_date",
				Direction: model.OrderDirectionAsc,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	a.Log.Trace("[END] - ", method)

	return data, total, nil
}

func (a *PayrollUseCase) CreatePeriod(
	ctx context.Context,
	request *model.CreatePayrollPeriodRequest,
	auth *model.Auth,
) error {
	method := "PayrollUseCase.CreatePeriod"
	a.Log.Trace("[BEGIN] - ", method)
	a.Log.Debug("request - ", method, request)

	db := a.DB.WithContext(ctx)

	startDate, err := time.Parse(time.DateOnly, request.StartDate)
	if err != nil {
		return fmt.Errorf("payroll-period/invalid-start-date")
	}
	endDate, err := time.Parse(time.DateOnly, request.EndDate)
	if err != nil {
		return fmt.Errorf("payroll-period/invalid-end-date")
	}

	payrollPeriod := entity.NewPayrollPeriod(&entity.CreatePayrollPeriodProps{
		StartDate: startDate,
		EndDate:   endDate,
		CreatedBy: auth.ID,
	})

	if !payrollPeriod.IsValidDateRange() {
		return fmt.Errorf("payroll-period/invalid-date-range")
	}

	isExist, err := a.payrollPeriod.IsExist(db, startDate, endDate)
	if err != nil {
		panic(err)
	} else if isExist {
		return fmt.Errorf("payroll/period-already-exists")
	}

	isOverlapping, err := a.payrollPeriod.IsOverlapping(db, startDate, endDate)
	if err != nil {
		panic(err)
	}

	if isOverlapping {
		return fmt.Errorf("payroll/period-overlapping")
	}

	if err := a.payrollPeriod.Create(db, payrollPeriod); err != nil {
		panic(err)
	}

	a.Log.Trace("[END] - ", method)

	return nil
}
