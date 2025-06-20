package middleware

import (
	"payslip-generator-service/config"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func SetupMiddleware(app *fiber.App, config *config.Config, logger *logrus.Logger) {
	app.Use(SetupHelmetMiddleware())
	app.Use(SetupRecoverMiddleware())
	app.Use(SetupCompressionMiddleware())
	app.Use(SetupCorsMiddleware(config))
	app.Use(SetupRequestIDMiddleware())
	app.Use(SetupRateLimiterMiddleware(config))
	app.Use(SetupCookieMiddleware(config))
	app.Use(SetupProbesMiddleware())
}

func SetupExceptionMiddleware(app *fiber.App) {
	app.Use(SetupNotFoundMiddleware())
}
