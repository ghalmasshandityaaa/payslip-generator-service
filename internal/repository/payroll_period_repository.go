package repository

import (
	"payslip-generator-service/internal/entity"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PayrollPeriodRepository struct {
	Repository[entity.PayrollPeriod]
	Log *logrus.Logger
}

func NewPayrollPeriodRepository(log *logrus.Logger) *PayrollPeriodRepository {
	return &PayrollPeriodRepository{
		Log: log,
	}
}

func (a *PayrollPeriodRepository) IsExist(db *gorm.DB, startDate, endDate time.Time) (bool, error) {
	var exists bool
	err := db.Model(&entity.PayrollPeriod{}).
		Select("1").
		Where("start_date = ? AND end_date = ?", startDate, endDate).
		Limit(1).
		Scan(&exists).Error

	return exists, err
}

func (a *PayrollPeriodRepository) IsOverlapping(db *gorm.DB, startDate, endDate time.Time) (bool, error) {
	var exists bool
	err := db.Model(&entity.PayrollPeriod{}).
		Select("1").
		Where("start_date <= ? AND end_date >= ?", endDate, startDate).
		Limit(1).
		Scan(&exists).Error

	return exists, err
}
