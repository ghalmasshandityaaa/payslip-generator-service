package route

func (a *Route) SetupReimbursementRoute() {
	a.Log.Info("setting up reimbursement routes")

	a.App.Post("/v1/reimbursement", a.AuthMiddleware, a.ReimbursementHandler.Create)
	a.Log.Info("mapped {/v1/reimbursement, POST} route")
}
