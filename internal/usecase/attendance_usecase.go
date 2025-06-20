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

type AttendanceUseCase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	AttendanceRepository *repository.AttendanceRepository
}

func NewAttendanceUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	employeeRepository *repository.AttendanceRepository,
) *AttendanceUseCase {
	return &AttendanceUseCase{
		DB:                   db,
		Log:                  logger,
		AttendanceRepository: employeeRepository,
	}
}

func (a *AttendanceUseCase) Create(
	ctx context.Context,
	request *model.CreateAttendanceRequest,
	auth *model.Auth,
) error {
	method := "AttendanceUseCase.Create"
	a.Log.Trace("[BEGIN] - ", method)
	a.Log.Debug("request - ", method, request)

	db := a.DB.WithContext(ctx)

	startTime, err := time.Parse(time.RFC3339, request.StartTime)
	if err != nil {
		return fmt.Errorf("attendance/invalid-start-time")
	}

	endTime, err := time.Parse(time.RFC3339, request.EndTime)
	if err != nil {
		return fmt.Errorf("attendance/invalid-end-time")
	}

	attendance := entity.NewAttendance(&entity.CreateAttendanceProps{
		StartTime: startTime,
		EndTime:   endTime,
		CreatedBy: auth.ID,
	})

	if !attendance.IsSameDay() {
		return fmt.Errorf("attendance/must-same-day")
	} else if !attendance.IsWeekday() {
		return fmt.Errorf("attendance/must-weekday")
	} else if !attendance.IsToday() {
		return fmt.Errorf("attendance/must-today")
	} else if !attendance.IsEndTimeGreaterThanStartTime() {
		return fmt.Errorf("attendance/invalid-time-order")
	}

	a.Log.Debug("attendance - ", method, attendance)

	if err := a.AttendanceRepository.Create(db, attendance); err != nil {
		return err
	}

	a.Log.Trace("[END] - ", method)

	return nil
}
