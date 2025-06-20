package route

func (a *Route) SetupAuthRoute() {
	a.Log.Info("setting up auth routes")

	a.App.Post("/v1/auth/sign-in", a.AuthHandler.SignIn)
	a.Log.Info("mapped {/v1/auth/sign-in, POST} route")
}
