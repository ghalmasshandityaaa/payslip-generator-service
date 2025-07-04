package repository

import (
	"payslip-generator-service/internal/entity"
	ulid "payslip-generator-service/pkg/database/gorm"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OvertimeRepository struct {
	Repository[entity.Overtime]
	Log *logrus.Logger
}

func NewOvertimeRepository(log *logrus.Logger) *OvertimeRepository {
	return &OvertimeRepository{
		Log: log,
	}
}

func (a *OvertimeRepository) FindByDate(db *gorm.DB, date time.Time) (*entity.Overtime, error) {
	var overtime entity.Overtime
	if err := db.Where("date = ?", date.Format(time.DateOnly)).First(&overtime).Error; err != nil {
		return nil, err
	}
	return &overtime, nil
}

func (a *OvertimeRepository) FindByPeriod(db *gorm.DB, employeeID ulid.ULID, startDate, endDate time.Time) ([]entity.Overtime, error) {
	var overtimes []entity.Overtime

	err := db.Debug().
		Where("date BETWEEN ? AND ?", startDate.Format(time.DateOnly), endDate.Format(time.DateOnly)).
		Where("created_by = ?", employeeID).
		Find(&overtimes).Error

	if err != nil {
		return nil, err
	}

	return overtimes, nil
}
