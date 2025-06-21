package repository

import (
	"payslip-generator-service/internal/entity"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AttendanceRepository struct {
	Repository[entity.Attendance]
	Log *logrus.Logger
}

func NewAttendanceRepository(log *logrus.Logger) *AttendanceRepository {
	return &AttendanceRepository{
		Log: log,
	}
}

func (a *AttendanceRepository) FindByDate(db *gorm.DB, date time.Time) (*entity.Attendance, error) {
	var attendance entity.Attendance
	if err := db.Where("DATE(start_time) = ?", date.Format(time.DateOnly)).First(&attendance).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (a *AttendanceRepository) FindByPeriod(db *gorm.DB, startDate, endDate time.Time) ([]entity.Attendance, error) {
	var attendances []entity.Attendance

	err := db.Debug().
		Where("DATE(start_time) BETWEEN ? AND ?", startDate.Format(time.DateOnly), endDate.Format(time.DateOnly)).
		Where("DATE(end_time) BETWEEN ? AND ?", startDate.Format(time.DateOnly), endDate.Format(time.DateOnly)).
		Where("DATE(created_at) BETWEEN ? AND ?", startDate.Format(time.DateOnly), endDate.Format(time.DateOnly)).
		Find(&attendances).Error

	if err != nil {
		return nil, err
	}

	return attendances, nil
}
