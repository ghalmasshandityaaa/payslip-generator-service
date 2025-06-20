package usecase

import (
	"context"
	"payslip-generator-service/internal/entity"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ReimbursementUseCase struct {
	DB                      *gorm.DB
	Log                     *logrus.Logger
	ReimbursementRepository *repository.ReimbursementRepository
}

func NewReimbursementUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	employeeRepository *repository.ReimbursementRepository,
) *ReimbursementUseCase {
	return &ReimbursementUseCase{
		DB:                      db,
		Log:                     logger,
		ReimbursementRepository: employeeRepository,
	}
}

func (a *ReimbursementUseCase) Create(
	ctx context.Context,
	request *model.CreateReimbursementRequest,
	auth *model.Auth,
) (*model.CreateReimbursementResponse, error) {
	method := "ReimbursementUseCase.Create"
	a.Log.Trace("[BEGIN] - ", method)
	a.Log.Debug("request - ", method, request)

	db := a.DB.WithContext(ctx)

	reimbursement := entity.NewReimbursement(&entity.CreateReimbursementProps{
		Amount:      request.Amount,
		Description: request.Description,
		CreatedBy:   auth.ID,
	})

	a.Log.Debug("reimbursement - ", method, reimbursement)

	if err := a.ReimbursementRepository.Create(db, reimbursement); err != nil {
		return nil, err
	}

	a.Log.Trace("[END] - ", method)

	return nil, nil
}
