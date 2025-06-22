package handler

import (
	"payslip-generator-service/config"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/usecase"
	"payslip-generator-service/pkg/logger"
	"payslip-generator-service/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Log       *logger.ContextLogger
	Config    *config.Config
	UseCase   *usecase.AuthUseCase
	Validator *validator.Validator
}

func NewAuthHandler(
	useCase *usecase.AuthUseCase,
	log *logger.ContextLogger,
	config *config.Config,
	validator *validator.Validator,
) *AuthHandler {
	return &AuthHandler{
		Log:       log,
		UseCase:   useCase,
		Config:    config,
		Validator: validator,
	}
}

// SignIn authenticates a user and returns access and refresh tokens
// @Summary Sign in user
// @Description Authenticate user with username and password to get access and refresh tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.SignInRequest true "Sign in credentials"
// @Router /auth/sign-in [post]
func (h *AuthHandler) SignIn(ctx *fiber.Ctx) error {
	method := "AuthHandler.SignIn"
	h.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")

	request := new(model.SignInRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}

	errValidation := h.Validator.ValidateStruct(request)
	if errValidation != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[*model.SignInResponse]{
			Ok:     false,
			Errors: errValidation,
		})
	}

	// Create context with request_id
	requestCtx := ctx.UserContext()
	accessToken, refreshToken, err := h.UseCase.SignIn(requestCtx, request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: err.Error(),
		})
	}

	h.Log.WithContext(ctx).WithField("method", method).Trace("[END]")

	return ctx.JSON(model.WebResponse[*model.SignInResponse]{
		Ok: true,
		Data: &model.SignInResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}
