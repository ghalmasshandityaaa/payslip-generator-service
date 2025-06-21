package app

import (
	"payslip-generator-service/config"
	"payslip-generator-service/internal/handler"
	"payslip-generator-service/internal/middleware"
	"payslip-generator-service/internal/repository"
	"payslip-generator-service/internal/route"
	"payslip-generator-service/internal/usecase"
	"payslip-generator-service/internal/utils"
	"payslip-generator-service/pkg/logger"
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
	// init context logger
	contextLogger := logger.NewContextLogger(config.Log)

	// init utils
	jwtUtil := utils.NewJwtUtil(config.Config)

	// init repositories
	userRepository := repository.NewEmployeeRepository(config.Log)
	reimbursementRepository := repository.NewReimbursementRepository(config.Log)
	attendanceRepository := repository.NewAttendanceRepository(config.Log)
	overtimeRepository := repository.NewOvertimeRepository(config.Log)
	payrollRepository := repository.NewPayrollPeriodRepository(config.Log)

	// init use cases
	authUseCase := usecase.NewAuthUseCase(config.DB, contextLogger, config.Config, jwtUtil, userRepository)
	employeeUseCase := usecase.NewEmployeeUseCase(config.DB, contextLogger, userRepository)
	reimbursementUseCase := usecase.NewReimbursementUseCase(config.DB, contextLogger, reimbursementRepository)
	attendanceUseCase := usecase.NewAttendanceUseCase(config.DB, contextLogger, attendanceRepository)
	overtimeUseCase := usecase.NewOvertimeUseCase(config.DB, contextLogger, overtimeRepository, attendanceRepository)
	payrollUseCase := usecase.NewPayrollUseCase(
		config.DB, contextLogger,
		payrollRepository,
		attendanceUseCase,
		overtimeUseCase,
		reimbursementUseCase,
		employeeUseCase,
	)

	// init handlers
	authHandler := handler.NewAuthHandler(authUseCase, contextLogger, config.Config, config.Validator)
	reimbursementHandler := handler.NewReimbursementHandler(reimbursementUseCase, contextLogger, config.Validator)
	attendanceHandler := handler.NewAttendanceHandler(attendanceUseCase, contextLogger, config.Validator)
	overtimeHandler := handler.NewOvertimeHandler(overtimeUseCase, contextLogger, config.Validator)
	payrollHandler := handler.NewPayrollHandler(payrollUseCase, contextLogger, config.Validator)

	// init middleware
	authMiddleware := middleware.NewAuthMiddleware(employeeUseCase, jwtUtil)

	// init routes
	appRoute := route.NewRoute(
		config.App, config.Log, authMiddleware,
		authHandler,
		reimbursementHandler,
		attendanceHandler,
		overtimeHandler,
		payrollHandler,
	)

	// setup routes
	appRoute.Setup()
}
