package handler

import (
	"context"
	"payslip-generator-service/internal/middleware"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/usecase"
	"payslip-generator-service/pkg/logger"
	"payslip-generator-service/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type OvertimeHandler struct {
	Log       *logger.ContextLogger
	UseCase   *usecase.OvertimeUseCase
	Validator *validator.Validator
}

func NewOvertimeHandler(
	useCase *usecase.OvertimeUseCase,
	log *logger.ContextLogger,
	validator *validator.Validator,
) *OvertimeHandler {
	return &OvertimeHandler{
		Log:       log,
		UseCase:   useCase,
		Validator: validator,
	}
}

// Create creates a new overtime record for the authenticated employee
// @Summary Create overtime record
// @Description Create a new overtime record with date and total hours for the authenticated employee
// @Tags Overtime
// @Accept json
// @Produce json
// @Security bearer
// @Param request body model.CreateOvertimeRequest true "Overtime details"
// @Router /overtime [post]
func (h *OvertimeHandler) Create(ctx *fiber.Ctx) error {
	method := "OvertimeHandler.Create"
	h.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")

	auth := middleware.GetAuth(ctx)
	request := new(model.CreateOvertimeRequest)
	if err := ctx.BodyParser(request); err != nil {
		h.Log.WithContext(ctx).Error("failed parse body: ", err.Error())
		return fiber.ErrBadRequest
	}

	errValidation := h.Validator.ValidateStruct(request)
	if errValidation != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: errValidation,
		})
	}

	// Create context with request_id
	requestCtx := context.WithValue(ctx.UserContext(), "request_id", logger.GetRequestID(ctx))
	err := h.UseCase.Create(requestCtx, request, auth)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: err.Error(),
		})
	}

	h.Log.WithContext(ctx).WithField("method", method).Trace("[END]")

	return ctx.JSON(model.WebResponse[*model.CreateOvertimeResponse]{
		Ok: true,
	})
}
