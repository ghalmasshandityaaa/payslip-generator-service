package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

func SetupProbesMiddleware() fiber.Handler {
	//TODO: implement complex probes configuration here
	return healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		ReadinessEndpoint: "/ready",
	})
}
