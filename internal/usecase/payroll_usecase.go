package usecase

import (
	"context"
	"errors"
	"fmt"
	"payslip-generator-service/internal/entity"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/repository"
	ulid "payslip-generator-service/pkg/database/gorm"
	"time"

	v2 "github.com/oklog/ulid/v2"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PayrollUseCase struct {
	DB                      *gorm.DB
	Log                     *logrus.Logger
	payrollPeriodRepository *repository.PayrollPeriodRepository
}

func NewPayrollUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	payrollPeriodRepository *repository.PayrollPeriodRepository,
) *PayrollUseCase {
	return &PayrollUseCase{
		DB:                      db,
		Log:                     logger,
		payrollPeriodRepository: payrollPeriodRepository,
	}
}

func (a *PayrollUseCase) ListPeriod(ctx context.Context, request *model.ListPayrollPeriodRequest) ([]entity.PayrollPeriod, int64, error) {
	method := "PayrollUseCase.ListPeriod"
	a.Log.Trace("[BEGIN] - ", method)
	a.Log.Debug("request - ", method, request)

	db := a.DB.WithContext(ctx)
	data, total, err := a.payrollPeriodRepository.FindAllWithPagination(db, &model.PaginationOptions{
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

	isExist, err := a.payrollPeriodRepository.IsExist(db, startDate, endDate)
	if err != nil {
		panic(err)
	} else if isExist {
		return fmt.Errorf("payroll/period-already-exists")
	}

	isOverlapping, err := a.payrollPeriodRepository.IsOverlapping(db, startDate, endDate)
	if err != nil {
		panic(err)
	}

	if isOverlapping {
		return fmt.Errorf("payroll/period-overlapping")
	}

	if err := a.payrollPeriodRepository.Create(db, payrollPeriod); err != nil {
		panic(err)
	}

	a.Log.Trace("[END] - ", method)

	return nil
}

func (a *PayrollUseCase) ProcessPayroll(ctx context.Context, request *model.ProcessPayrollRequest, auth *model.Auth) error {
	method := "PayrollUseCase.ProcessPayroll"
	a.Log.Trace("[BEGIN] - ", method)
	a.Log.Debug("request - ", method, request)

	db := a.DB.WithContext(ctx)

	payrollPeriod := new(entity.PayrollPeriod)
	err := a.payrollPeriodRepository.FindById(db, payrollPeriod, ulid.ULID(v2.MustParse(request.PeriodID)))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("payroll/period-not-found")
		}
		panic(err)
	}

	if payrollPeriod.IsProcessed() {
		return fmt.Errorf("payroll/already-processed")
	}

	payrollPeriod.Process(auth.ID)
	if err := a.payrollPeriodRepository.Update(db, payrollPeriod); err != nil {
		panic(err)
	}

	a.Log.Trace("[END] - ", method)

	return nil
}
