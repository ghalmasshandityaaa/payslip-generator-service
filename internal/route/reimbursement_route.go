package route

import (
	"payslip-generator-service/internal/middleware"
	"payslip-generator-service/internal/model"
)

func (a *Route) SetupReimbursementRoute() {
	a.Log.Info("setting up reimbursement routes")

	a.App.Post("/v1/reimbursement", a.AuthMiddleware, middleware.RoleMiddleware(model.RoleEmployee), a.ReimbursementHandler.Create)
	a.Log.Info("mapped {/v1/reimbursement, POST} route")
}
