package usecase

import (
	"context"
	"payslip-generator-service/config"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/repository"
	"payslip-generator-service/internal/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Config             *config.Config
	JwtUtil            *utils.JwtUtil
	EmployeeRepository *repository.EmployeeRepository
}

func NewAuthUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	config *config.Config,
	jwtUtil *utils.JwtUtil,
	employeeRepository *repository.EmployeeRepository,
) *AuthUseCase {
	return &AuthUseCase{
		DB:                 db,
		Log:                logger,
		Config:             config,
		JwtUtil:            jwtUtil,
		EmployeeRepository: employeeRepository,
	}
}

func (a *AuthUseCase) SignIn(ctx context.Context, request *model.SignInRequest) (string, string, error) {
	return "accessToken", "refreshToken", nil
}
