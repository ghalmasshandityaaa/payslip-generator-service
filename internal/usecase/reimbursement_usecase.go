package usecase

import (
	"context"
	"payslip-generator-service/internal/entity"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/repository"
	"payslip-generator-service/pkg/logger"
	"time"

	ulid "payslip-generator-service/pkg/database/gorm"

	"gorm.io/gorm"
)

type ReimbursementUseCase struct {
	DB                      *gorm.DB
	Log                     *logger.ContextLogger
	ReimbursementRepository *repository.ReimbursementRepository
}

func NewReimbursementUseCase(
	db *gorm.DB,
	log *logger.ContextLogger,
	reimbursementRepository *repository.ReimbursementRepository,
) *ReimbursementUseCase {
	return &ReimbursementUseCase{
		DB:                      db,
		Log:                     log,
		ReimbursementRepository: reimbursementRepository,
	}
}

func (a *ReimbursementUseCase) Create(
	ctx context.Context,
	request *model.CreateReimbursementRequest,
	auth *model.Auth,
) error {
	method := "ReimbursementUseCase.Create"
	a.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")
	a.Log.WithContext(ctx).WithField("method", method).WithField("request", request).Debug("request")

	db := a.DB.WithContext(ctx)

	reimbursement := entity.NewReimbursement(&entity.CreateReimbursementProps{
		Amount:      request.Amount,
		Description: request.Description,
		CreatedBy:   auth.ID,
	})

	a.Log.WithContext(ctx).Debug("reimbursement - ", method, reimbursement)

	if err := a.ReimbursementRepository.Create(db, reimbursement); err != nil {
		panic(err)
	}

	a.Log.WithContext(ctx).WithField("method", method).Trace("[END]")

	return nil
}

func (a *ReimbursementUseCase) ListByPeriod(
	ctx context.Context,
	employeeID ulid.ULID,
	startDate time.Time,
	endDate time.Time,
) ([]entity.Reimbursement, error) {
	method := "ReimbursementUseCase.ListByPeriod"
	a.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")
	a.Log.WithContext(ctx).WithField("method", method).WithField("request", startDate).WithField("request", endDate).Debug("request")

	db := a.DB.WithContext(ctx)

	reimbursements, err := a.ReimbursementRepository.FindByPeriod(db, employeeID, startDate, endDate)
	if err != nil {
		panic(err)
	}

	a.Log.WithContext(ctx).WithField("method", method).Trace("[END]")

	return reimbursements, nil
}
