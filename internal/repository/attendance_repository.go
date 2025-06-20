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
