package handler

import (
	"math"
	"payslip-generator-service/internal/entity"
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

func (h *PayrollHandler) ListPeriod(ctx *fiber.Ctx) error {
	method := "PayrollHandler.ListPeriod"
	h.Log.Trace("[BEGIN] - ", method)

	request := &model.ListPayrollPeriodRequest{
		Page:     ctx.QueryInt("page", 1),
		PageSize: ctx.QueryInt("size", 10),
	}

	data, total, err := h.UseCase.ListPeriod(ctx.UserContext(), request)
	if err != nil {
		return ctx.JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: err.Error(),
		})
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		PageSize:  request.PageSize,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.PageSize))),
	}

	h.Log.Trace("[END] - ", method)
	return ctx.JSON(model.WebResponseWithData[[]entity.PayrollPeriod]{
		Ok:     true,
		Data:   data,
		Paging: paging,
	})
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

func (h *PayrollHandler) ProcessPayroll(ctx *fiber.Ctx) error {
	method := "PayrollHandler.ProcessPayroll"
	h.Log.Trace("[BEGIN] - ", method)

	auth := middleware.GetAuth(ctx)
	request := new(model.ProcessPayrollRequest)
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

	err := h.UseCase.ProcessPayroll(ctx.UserContext(), request, auth)
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
