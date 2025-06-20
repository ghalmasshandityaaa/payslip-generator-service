package route

import (
	"payslip-generator-service/internal/middleware"
	"payslip-generator-service/internal/model"
)

func (a *Route) SetupOvertimeRoute() {
	a.Log.Info("setting up overtime routes")

	a.App.Post("/v1/overtime", a.AuthMiddleware, middleware.RoleMiddleware(model.RoleEmployee), a.OvertimeHandler.Create)
	a.Log.Info("mapped {/v1/overtime, POST} route")
}
