package route

import (
	"payslip-generator-service/internal/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Route struct {
	App                  *fiber.App
	Log                  *logrus.Logger
	AuthMiddleware       fiber.Handler
	AuthHandler          *handler.AuthHandler
	ReimbursementHandler *handler.ReimbursementHandler
	AttendanceHandler    *handler.AttendanceHandler
}

func NewRoute(
	app *fiber.App,
	logger *logrus.Logger,
	authMiddleware fiber.Handler,
	authHandler *handler.AuthHandler,
	reimbursementHandler *handler.ReimbursementHandler,
	attendanceHandler *handler.AttendanceHandler,
) *Route {
	return &Route{
		App:                  app,
		Log:                  logger,
		AuthMiddleware:       authMiddleware,
		AuthHandler:          authHandler,
		ReimbursementHandler: reimbursementHandler,
		AttendanceHandler:    attendanceHandler,
	}
}

func (a *Route) Setup() {
	a.Log.Info("setting up routes")

	a.SetupAuthRoute()
	a.SetupReimbursementRoute()
	a.SetupAttendanceRoute()
}
