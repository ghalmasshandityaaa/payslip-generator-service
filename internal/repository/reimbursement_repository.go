package repository

import (
	"payslip-generator-service/internal/entity"

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

// GetByID retrieves a reimbursement by ID
func (r *ReimbursementRepository) GetByID(db *gorm.DB, reimbursement *entity.Reimbursement, id string) error {
	return db.Debug().Where("id = ?", id).Take(reimbursement).Error
}

// GetByCreatedBy retrieves reimbursements created by a specific employee
func (r *ReimbursementRepository) GetByCreatedBy(db *gorm.DB, reimbursements *[]entity.Reimbursement, createdBy string) error {
	return db.Debug().Where("created_by = ?", createdBy).Find(reimbursements).Error
}

// GetByStatus retrieves reimbursements by status
func (r *ReimbursementRepository) GetByStatus(db *gorm.DB, reimbursements *[]entity.Reimbursement, status entity.ReimbursementStatus) error {
	return db.Debug().Where("status = ?", status).Find(reimbursements).Error
}
