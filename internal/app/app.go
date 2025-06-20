package app

import (
	"payslip-generator-service/config"
	"payslip-generator-service/internal/handler"
	"payslip-generator-service/internal/middleware"
	"payslip-generator-service/internal/repository"
	"payslip-generator-service/internal/route"
	"payslip-generator-service/internal/usecase"
	"payslip-generator-service/internal/utils"
	"payslip-generator-service/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	App       *fiber.App
	Log       *logrus.Logger
	Config    *config.Config
	DB        *gorm.DB
	Validator *validator.Validator
}

func Bootstrap(config *BootstrapConfig) {
	// init utils
	jwtUtil := utils.NewJwtUtil(config.Config)

	// init repositories
	userRepository := repository.NewEmployeeRepository(config.Log)
	reimbursementRepository := repository.NewReimbursementRepository(config.Log)
	attendanceRepository := repository.NewAttendanceRepository(config.Log)
	overtimeRepository := repository.NewOvertimeRepository(config.Log)

	// init use cases
	authUseCase := usecase.NewAuthUseCase(config.DB, config.Log, config.Config, jwtUtil, userRepository)
	employeeUseCase := usecase.NewEmployeeUseCase(config.DB, config.Log, userRepository)
	reimbursementUseCase := usecase.NewReimbursementUseCase(config.DB, config.Log, reimbursementRepository)
	attendanceUseCase := usecase.NewAttendanceUseCase(config.DB, config.Log, attendanceRepository)
	overtimeUseCase := usecase.NewOvertimeUseCase(config.DB, config.Log, overtimeRepository, attendanceRepository)

	// init handlers
	authHandler := handler.NewAuthHandler(authUseCase, config.Log, config.Config, config.Validator)
	reimbursementHandler := handler.NewReimbursementHandler(reimbursementUseCase, config.Log, config.Validator)
	attendanceHandler := handler.NewAttendanceHandler(attendanceUseCase, config.Log, config.Validator)
	overtimeHandler := handler.NewOvertimeHandler(overtimeUseCase, config.Log, config.Validator)

	// init middleware
	authMiddleware := middleware.NewAuthMiddleware(employeeUseCase, jwtUtil)

	// init routes
	appRoute := route.NewRoute(config.App, config.Log, authMiddleware, authHandler, reimbursementHandler, attendanceHandler, overtimeHandler)

	// setup routes
	appRoute.Setup()
}
