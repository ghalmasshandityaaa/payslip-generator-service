package usecase

import (
	"context"
	"fmt"
	"payslip-generator-service/internal/entity"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/repository"
	"time"

	ulid "payslip-generator-service/pkg/database/gorm"

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
	attendanceRepository *repository.AttendanceRepository,
) *AttendanceUseCase {
	return &AttendanceUseCase{
		DB:                   db,
		Log:                  logger,
		AttendanceRepository: attendanceRepository,
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

	todayAttendance, err := a.AttendanceRepository.FindByDate(db, attendance.StartTime)
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	if todayAttendance != nil {
		return fmt.Errorf("attendance/already-exists")
	}

	if err := a.AttendanceRepository.Create(db, attendance); err != nil {
		panic(err)
	}

	a.Log.Trace("[END] - ", method)

	return nil
}

func (a *AttendanceUseCase) ListByPeriod(
	ctx context.Context,
	employeeID ulid.ULID,
	startDate time.Time,
	endDate time.Time,
) ([]entity.Attendance, error) {
	method := "AttendanceUseCase.ListByPeriod"
	a.Log.Trace("[BEGIN] - ", method)
	a.Log.Debug("request - ", method, startDate, endDate)

	db := a.DB.WithContext(ctx)

	attendances, err := a.AttendanceRepository.FindByPeriod(db, employeeID, startDate, endDate)
	if err != nil {
		panic(err)
	}

	a.Log.Trace("[END] - ", method)

	return attendances, nil
}
