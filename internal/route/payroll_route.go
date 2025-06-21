package route

import (
	"payslip-generator-service/internal/middleware"
	"payslip-generator-service/internal/model"
)

func (a *Route) SetupPayrollRoute() {
	a.Log.Info("setting up payroll routes")

	a.App.Get("/v1/payroll/period", a.AuthMiddleware, a.PayrollHandler.ListPeriod)
	a.Log.Info("mapped {/v1/payroll/period, GET} route")

	a.App.Post("/v1/payroll/period", a.AuthMiddleware, middleware.RoleMiddleware(model.RoleAdmin), a.PayrollHandler.CreatePeriod)
	a.Log.Info("mapped {/v1/payroll/period, POST} route")

	a.App.Post("/v1/payroll/process", a.AuthMiddleware, middleware.RoleMiddleware(model.RoleAdmin), a.PayrollHandler.ProcessPayroll)
	a.Log.Info("mapped {/v1/payroll/process, POST} route")

	a.App.Get("/v1/payroll/payslip", a.AuthMiddleware, middleware.RoleMiddleware(model.RoleEmployee), a.PayrollHandler.GetPayslip)
	a.Log.Info("mapped {/v1/payroll/payslip, GET} route")

	a.App.Get("/v1/payroll/payslip/report", a.AuthMiddleware, middleware.RoleMiddleware(model.RoleAdmin), a.PayrollHandler.GetPayslipReport)
	a.Log.Info("mapped {/v1/payroll/payslip/report, GET} route")
}
