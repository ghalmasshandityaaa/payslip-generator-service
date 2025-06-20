package handler

import (
	"payslip-generator-service/internal/middleware"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/usecase"
	"payslip-generator-service/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type PayrollHandler struct {
	Log       *logrus.Logger
	UseCase   *usecase.PayrollUseCase
	Validator *validator.Validator
}

func NewPayrollHandler(
	useCase *usecase.PayrollUseCase,
	logger *logrus.Logger,
	validator *validator.Validator,
) *PayrollHandler {
	return &PayrollHandler{
		Log:       logger,
		UseCase:   useCase,
		Validator: validator,
	}
}

func (h *PayrollHandler) CreatePeriod(ctx *fiber.Ctx) error {
	method := "PayrollHandler.Create"
	h.Log.Trace("[BEGIN] - ", method)

	auth := middleware.GetAuth(ctx)
	request := new(model.CreatePayrollPeriodRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}

	errValidation := h.Validator.ValidateStruct(request)
	if errValidation != nil {
		return ctx.JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: errValidation,
		})
	}

	err := h.UseCase.CreatePeriod(ctx.UserContext(), request, auth)
	if err != nil {
		return ctx.JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: err.Error(),
		})
	}

	h.Log.Trace("[END] - ", method)
	return ctx.JSON(model.WebResponse[*model.CreatePayrollPeriodResponse]{
		Ok: true,
	})
}
