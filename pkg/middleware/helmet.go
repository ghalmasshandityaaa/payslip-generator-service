package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

func SetupHelmetMiddleware() fiber.Handler {
	return helmet.New(helmet.ConfigDefault)
}
