package middleware

import (
	"context"
	"payslip-generator-service/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/oklog/ulid/v2"
)

func SetupRequestIDMiddleware() fiber.Handler {
	return requestid.New(requestid.Config{
		Header: "X-Request-ID",
		Generator: func() string {
			return ulid.Make().String()
		},
		ContextKey: "request_id",
		Next: func(c *fiber.Ctx) bool {
			// Store additional request information in context
			requestCtx := context.WithValue(c.UserContext(), "request_id", c.Locals("request_id"))
			requestCtx = context.WithValue(requestCtx, "ip_address", logger.GetIPAddress(c))
			requestCtx = context.WithValue(requestCtx, "user_agent", logger.GetUserAgent(c))

			// Update the context in Fiber
			c.SetUserContext(requestCtx)
			return false
		},
	})
}
