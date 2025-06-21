package usecase

import (
	"context"
	"errors"
	"fmt"
	"payslip-generator-service/internal/entity"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/repository"
	"payslip-generator-service/internal/vm"
	ulid "payslip-generator-service/pkg/database/gorm"
	"time"

	v2 "github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type PayrollUseCase struct {
	DB                      *gorm.DB
	Log                     *logrus.Logger
	payrollPeriodRepository *repository.PayrollPeriodRepository
	attendanceUseCase       *AttendanceUseCase
	overtimeUseCase         *OvertimeUseCase
	reimbursementUseCase    *ReimbursementUseCase
	employeeUseCase         *EmployeeUseCase
}

func NewPayrollUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	payrollPeriodRepository *repository.PayrollPeriodRepository,
	attendanceUseCase *AttendanceUseCase,
	overtimeUseCase *OvertimeUseCase,
	reimbursementUseCase *ReimbursementUseCase,
	employeeUseCase *EmployeeUseCase,
) *PayrollUseCase {
	return &PayrollUseCase{
		DB:                      db,
		Log:                     logger,
		payrollPeriodRepository: payrollPeriodRepository,
		attendanceUseCase:       attendanceUseCase,
		overtimeUseCase:         overtimeUseCase,
		reimbursementUseCase:    reimbursementUseCase,
		employeeUseCase:         employeeUseCase,
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

func (a *PayrollUseCase) generatePayslip(
	ctx context.Context,
	params model.GeneratePayslipRequest,
) (*vm.Payslip, error) {
	method := "PayrollUseCase.generatePayslip"
	a.Log.Trace("[BEGIN] - ", method)
	a.Log.Debug("request - ", method, params)

	if params.Period.ProcessedAt == nil {
		return nil, fmt.Errorf("payroll/period-not-processed")
	}

	var (
		attendance    []entity.Attendance
		overtime      []entity.Overtime
		reimbursement []entity.Reimbursement
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				a.Log.Error("Panic in attendance goroutine:", r)
				if err, ok := r.(error); ok {
					returnErr = err
				} else {
					returnErr = fmt.Errorf("panic: %v", r)
				}
			}
		}()

		var err error
		attendance, err = a.attendanceUseCase.ListByPeriod(ctx, params.EmployeeID, params.Period.StartDate, params.Period.EndDate)
		return err
	})

	g.Go(func() (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				a.Log.Error("Panic in overtime goroutine:", r)
				if err, ok := r.(error); ok {
					returnErr = err
				} else {
					returnErr = fmt.Errorf("panic: %v", r)
				}
			}
		}()

		var err error
		overtime, err = a.overtimeUseCase.ListByPeriod(ctx, params.EmployeeID, params.Period.StartDate, params.Period.EndDate)
		return err
	})

	g.Go(func() (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				a.Log.Error("Panic in reimbursement goroutine:", r)
				if err, ok := r.(error); ok {
					returnErr = err
				} else {
					returnErr = fmt.Errorf("panic: %v", r)
				}
			}
		}()

		var err error
		reimbursement, err = a.reimbursementUseCase.ListByPeriod(ctx, params.EmployeeID, params.Period.StartDate, *params.Period.ProcessedAt)
		return err
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}

	a.Log.Info("Generate payslip for period: ", params.Period.StartDate, " to ", params.Period.EndDate)

	payslip := vm.NewPayslip(&vm.CreatePayslipProps{
		EmployeeID:    params.EmployeeID,
		Attendance:    attendance,
		Overtime:      overtime,
		Reimbursement: reimbursement,
		PayrollPeriod: params.Period,
		Salary:        params.Salary,
	})

	a.Log.Trace("[END] - ", method)
	return payslip, nil
}

func (a *PayrollUseCase) GetPayslip(ctx context.Context, request *model.GetPayslipRequest, auth *model.Auth) (*vm.Payslip, error) {
	method := "PayrollUseCase.GetPayslip"
	a.Log.Trace("[BEGIN] - ", method)
	a.Log.Debug("request - ", method, request)

	db := a.DB.WithContext(ctx)

	payrollPeriod := new(entity.PayrollPeriod)
	err := a.payrollPeriodRepository.FindById(db, payrollPeriod, ulid.ULID(v2.MustParse(request.PeriodID)))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("payroll/period-not-found")
		}
		panic(err)
	}

	if !payrollPeriod.IsProcessed() {
		return nil, fmt.Errorf("payroll/not-processed")
	}

	employee, err := a.employeeUseCase.GetById(ctx, auth.ID)
	if err != nil {
		panic(err)
	}

	payslip, err := a.generatePayslip(ctx, model.GeneratePayslipRequest{
		EmployeeID: employee.ID,
		Salary:     employee.Salary,
		Period:     *payrollPeriod,
	})
	if err != nil {
		panic(err)
	}

	a.Log.Trace("[END] - ", method)

	return payslip, nil
}

func (a *PayrollUseCase) GetPayslipReport(ctx context.Context, request *model.GetPayslipRequest, auth *model.Auth) (*vm.PayslipReport, error) {
	method := "PayrollUseCase.GetPayslip"
	a.Log.Trace("[BEGIN] - ", method)
	a.Log.Debug("request - ", method, request)

	db := a.DB.WithContext(ctx)

	payrollPeriod := new(entity.PayrollPeriod)
	err := a.payrollPeriodRepository.FindById(db, payrollPeriod, ulid.ULID(v2.MustParse(request.PeriodID)))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("payroll/period-not-found")
		}
		panic(err)
	}

	if !payrollPeriod.IsProcessed() {
		return nil, fmt.Errorf("payroll/not-processed")
	}

	employees, err := a.employeeUseCase.List(ctx)
	if err != nil {
		panic(err)
	}

	var payslips []vm.Payslip
	g, ctx := errgroup.WithContext(ctx)

	for _, employee := range employees {
		g.Go(func() (returnErr error) {
			defer func() {
				if r := recover(); r != nil {
					a.Log.Error("Panic in payslip goroutine:", r)
					if err, ok := r.(error); ok {
						returnErr = err
					} else {
						returnErr = fmt.Errorf("panic: %v", r)
					}
				}
			}()

			var err error
			payslip, err := a.generatePayslip(ctx, model.GeneratePayslipRequest{
				EmployeeID: employee.ID,
				Salary:     employee.Salary,
				Period:     *payrollPeriod,
			})
			payslips = append(payslips, *payslip)
			return err
		})
	}

	if err := g.Wait(); err != nil {
		panic(err)
	}

	payslipReport := vm.NewPayslipReport(&vm.CreatePayslipReportProps{
		Employees: employees,
		Payslips:  payslips,
	})

	return payslipReport, nil
}
