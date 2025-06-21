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

type OvertimeUseCase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	OvertimeRepository   *repository.OvertimeRepository
	AttendanceRepository *repository.AttendanceRepository
}

func NewOvertimeUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	overtimeRepository *repository.OvertimeRepository,
	attendanceRepository *repository.AttendanceRepository,
) *OvertimeUseCase {
	return &OvertimeUseCase{
		DB:                   db,
		Log:                  logger,
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
	a.Log.Trace("[BEGIN] - ", method)
	a.Log.Debug("request - ", method, request)

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

	a.Log.Trace("[END] - ", method)

	return nil
}
func (a *OvertimeUseCase) ListByPeriod(
	ctx context.Context,
	auth *model.Auth,
	startDate time.Time,
	endDate time.Time,
) ([]entity.Overtime, error) {
	method := "OvertimeUseCase.ListByPeriod"
	a.Log.Trace("[BEGIN] - ", method)
	a.Log.Debug("request - ", method, startDate, endDate)

	db := a.DB.WithContext(ctx)

	overtimes, err := a.OvertimeRepository.FindByPeriod(db, startDate, endDate)
	if err != nil {
		panic(err)
	}

	a.Log.Trace("[END] - ", method)

	return overtimes, nil
}
