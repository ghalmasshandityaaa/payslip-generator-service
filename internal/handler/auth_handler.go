package handler

import (
	"payslip-generator-service/config"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/usecase"
	"payslip-generator-service/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	Log       *logrus.Logger
	Config    *config.Config
	UseCase   *usecase.AuthUseCase
	Validator *validator.Validator
}

func NewAuthHandler(
	useCase *usecase.AuthUseCase,
	logger *logrus.Logger,
	config *config.Config,
	validator *validator.Validator,
) *AuthHandler {
	return &AuthHandler{
		Log:       logger,
		UseCase:   useCase,
		Config:    config,
		Validator: validator,
	}
}

func (h *AuthHandler) SignIn(ctx *fiber.Ctx) error {
	method := "AuthHandler.SignIn"
	h.Log.Trace(method, "BEGIN")

	request := new(model.SignInRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}

	errValidation := h.Validator.ValidateStruct(request)
	if errValidation != nil {
		return ctx.JSON(model.WebResponse[*model.SignInResponse]{
			Ok:     false,
			Errors: errValidation,
		})
	}

	accessToken, refreshToken, err := h.UseCase.SignIn(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	h.Log.Trace(method, "END")

	return ctx.JSON(model.WebResponse[*model.SignInResponse]{
		Ok: true,
		Data: &model.SignInResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}
