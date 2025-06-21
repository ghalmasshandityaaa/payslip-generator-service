package repository

import (
	"payslip-generator-service/internal/entity"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ReimbursementRepository struct {
	Repository[entity.Reimbursement]
	Log *logrus.Logger
}

func NewReimbursementRepository(log *logrus.Logger) *ReimbursementRepository {
	return &ReimbursementRepository{
		Log: log,
	}
}

func (a *ReimbursementRepository) FindByPeriod(db *gorm.DB, startDate, endDate time.Time) ([]entity.Reimbursement, error) {
	var reimbursements []entity.Reimbursement

	err := db.
		Debug().
		Where("DATE(created_at) BETWEEN ? AND ?", startDate.Format(time.DateOnly), endDate.Format(time.DateOnly)).
		Find(&reimbursements).Error

	if err != nil {
		return nil, err
	}

	return reimbursements, nil
}
