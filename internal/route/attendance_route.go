package route

import (
	"payslip-generator-service/internal/middleware"
	"payslip-generator-service/internal/model"
)

func (a *Route) SetupAttendanceRoute() {
	a.Log.Info("setting up attendance routes")

	a.App.Post("/v1/attendance", a.AuthMiddleware, middleware.RoleMiddleware(model.RoleEmployee), a.AttendanceHandler.Create)
	a.Log.Info("mapped {/v1/attendance, POST} route")
}
