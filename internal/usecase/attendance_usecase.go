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

type AttendanceUseCase struct {
	DB                   *gorm.DB
	Log                  *logger.ContextLogger
	AttendanceRepository *repository.AttendanceRepository
}

func NewAttendanceUseCase(
	db *gorm.DB,
	log *logger.ContextLogger,
	attendanceRepository *repository.AttendanceRepository,
) *AttendanceUseCase {
	return &AttendanceUseCase{
		DB:                   db,
		Log:                  log,
		AttendanceRepository: attendanceRepository,
	}
}

func (a *AttendanceUseCase) Create(
	ctx context.Context,
	request *model.CreateAttendanceRequest,
	auth *model.Auth,
) error {
	method := "AttendanceUseCase.Create"
	a.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")
	a.Log.WithContext(ctx).WithField("method", method).WithField("request", request).Debug("request")

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
	} else if !attendance.IsEndTimeGreaterThanStartTime() {
		return fmt.Errorf("attendance/invalid-time-order")
	} else if !attendance.IsToday() {
		return fmt.Errorf("attendance/must-today")
	}

	a.Log.WithContext(ctx).Debug("attendance - ", method, attendance)

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

	a.Log.WithContext(ctx).WithField("method", method).Trace("[END]")

	return nil
}

func (a *AttendanceUseCase) ListByPeriod(
	ctx context.Context,
	employeeID ulid.ULID,
	startDate time.Time,
	endDate time.Time,
) ([]entity.Attendance, error) {
	method := "AttendanceUseCase.ListByPeriod"
	a.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")
	a.Log.WithContext(ctx).WithField("method", method).WithField("request", startDate).WithField("request", endDate).Debug("request")

	db := a.DB.WithContext(ctx)

	attendances, err := a.AttendanceRepository.FindByPeriod(db, employeeID, startDate, endDate)
	if err != nil {
		panic(err)
	}

	a.Log.WithContext(ctx).WithField("method", method).Trace("[END]")

	return attendances, nil
}
