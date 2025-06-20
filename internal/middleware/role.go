package middleware

import (
	"payslip-generator-service/internal/model"

	"github.com/gofiber/fiber/v2"
)

func RoleMiddleware(role model.Role) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := GetAuth(ctx)
		if auth == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(model.WebResponse[any]{
				Ok:     false,
				Errors: "auth/unauthorized",
			})
		}

		if role == model.RoleAdmin {
			if !auth.IsAdmin {
				return ctx.Status(fiber.StatusForbidden).JSON(model.WebResponse[any]{
					Ok:     false,
					Errors: "auth/forbidden-access",
				})
			}
		} else {
			if auth.IsAdmin {
				return ctx.Status(fiber.StatusForbidden).JSON(model.WebResponse[any]{
					Ok:     false,
					Errors: "auth/forbidden-access",
				})
			}
		}

		return ctx.Next()
	}
}
