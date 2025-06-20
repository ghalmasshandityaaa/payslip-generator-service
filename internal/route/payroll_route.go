package route

import (
	"payslip-generator-service/internal/middleware"
	"payslip-generator-service/internal/model"
)

func (a *Route) SetupPayrollRoute() {
	a.Log.Info("setting up payroll routes")

	a.App.Post("/v1/payroll/period", a.AuthMiddleware, middleware.RoleMiddleware(model.RoleAdmin), a.PayrollHandler.CreatePeriod)
	a.Log.Info("mapped {/v1/payroll/period, POST} route")
}
