package repository

import (
	"payslip-generator-service/internal/entity"
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
