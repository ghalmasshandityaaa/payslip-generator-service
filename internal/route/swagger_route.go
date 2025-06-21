package route

import (
	"github.com/gofiber/swagger"
)

func (a *Route) SetupSwaggerRoute() {
	a.Log.Info("setting up swagger routes")

	a.App.Get("/swagger/*", swagger.HandlerDefault)
	a.Log.Info("mapped {/swagger/*, GET} route")
}
