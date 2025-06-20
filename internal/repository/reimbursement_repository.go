package repository

import (
	"payslip-generator-service/internal/entity"

	"github.com/sirupsen/logrus"
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
