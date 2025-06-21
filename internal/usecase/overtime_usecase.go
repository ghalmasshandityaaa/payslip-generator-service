package usecase

import (
	"context"
	"fmt"
	"payslip-generator-service/internal/entity"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/repository"
	"payslip-generator-service/pkg/logger"
	"time"

	ulid "payslip-generator-service/pkg/database/gorm"

	"gorm.io/gorm"
)

type OvertimeUseCase struct {
	DB                   *gorm.DB
	Log                  *logger.ContextLogger
	OvertimeRepository   *repository.OvertimeRepository
	AttendanceRepository *repository.AttendanceRepository
}

func NewOvertimeUseCase(
	db *gorm.DB,
	log *logger.ContextLogger,
	overtimeRepository *repository.OvertimeRepository,
	attendanceRepository *repository.AttendanceRepository,
) *OvertimeUseCase {
	return &OvertimeUseCase{
		DB:                   db,
		Log:                  log,
		OvertimeRepository:   overtimeRepository,
		AttendanceRepository: attendanceRepository,
	}
}

func (a *OvertimeUseCase) Create(
	ctx context.Context,
	request *model.CreateOvertimeRequest,
	auth *model.Auth,
) error {
	method := "OvertimeUseCase.Create"
	a.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")
	a.Log.WithContext(ctx).WithField("method", method).WithField("request", request).Debug("request")

	db := a.DB.WithContext(ctx)

	date, err := time.Parse(time.DateOnly, request.Date)
	if err != nil {
		return fmt.Errorf("overtime/invalid-date")
	}

	_, err = a.AttendanceRepository.FindByDate(db, date)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("attendance/not-found")
		}
		panic(err)
	}

	todayOvertime, err := a.OvertimeRepository.FindByDate(db, date)
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	if todayOvertime != nil {
		return fmt.Errorf("overtime/already-exists")
	}

	overtime := entity.NewOvertime(&entity.CreateOvertimeProps{
		Date:       date,
		TotalHours: request.TotalHours,
		CreatedBy:  auth.ID,
	})

	if !overtime.IsValidDuration() {
		return fmt.Errorf("overtime/invalid-duration")
	} else if !overtime.IsToday() {
		return fmt.Errorf("overtime/must-today")
	}

	if err := a.OvertimeRepository.Create(db, overtime); err != nil {
		panic(err)
	}

	a.Log.WithContext(ctx).WithField("method", method).Trace("[END]")

	return nil
}

func (a *OvertimeUseCase) ListByPeriod(
	ctx context.Context,
	employeeID ulid.ULID,
	startDate time.Time,
	endDate time.Time,
) ([]entity.Overtime, error) {
	method := "OvertimeUseCase.ListByPeriod"
	a.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")
	a.Log.WithContext(ctx).WithField("method", method).WithField("request", startDate).WithField("request", endDate).Debug("request")

	db := a.DB.WithContext(ctx)

	overtimes, err := a.OvertimeRepository.FindByPeriod(db, employeeID, startDate, endDate)
	if err != nil {
		panic(err)
	}

	a.Log.WithContext(ctx).WithField("method", method).Trace("[END]")

	return overtimes, nil
}
