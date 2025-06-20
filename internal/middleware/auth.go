package middleware

import (
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/usecase"
	"payslip-generator-service/internal/utils"
	"strings"

	"github.com/oklog/ulid/v2"

	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware(employeeUseCase *usecase.EmployeeUseCase, util *utils.JwtUtil) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := &model.VerifyAccountRequest{
			Token: ctx.Get("Authorization", "NOT_FOUND"),
		}

		token := request.Token
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		claim, err := util.ValidateToken(token)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(model.WebResponse[any]{
				Ok:     false,
				Errors: "auth/unauthorized",
			})
		}

		employee, err := employeeUseCase.GetById(ctx.UserContext(), ulid.MustParse(claim.ID))
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(model.WebResponse[any]{
				Ok:     false,
				Errors: "auth/unauthorized",
			})
		}

		ctx.Locals("auth", &model.Auth{
			ID:      employee.ID,
			IsAdmin: employee.IsAdmin,
		})
		return ctx.Next()
	}
}

func GetAuth(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
