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

type AttendanceHandler struct {
	Log       *logger.ContextLogger
	UseCase   *usecase.AttendanceUseCase
	Validator *validator.Validator
}

func NewAttendanceHandler(
	useCase *usecase.AttendanceUseCase,
	log *logger.ContextLogger,
	validator *validator.Validator,
) *AttendanceHandler {
	return &AttendanceHandler{
		Log:       log,
		UseCase:   useCase,
		Validator: validator,
	}
}

// Create creates a new attendance record for the authenticated employee
// @Summary Create attendance record
// @Description Create a new attendance record with start and end times for the authenticated employee
// @Tags Attendance
// @Accept json
// @Produce json
// @Security bearer
// @Param request body model.CreateAttendanceRequest true "Attendance details"
// @Router /attendance [post]
func (h *AttendanceHandler) Create(ctx *fiber.Ctx) error {
	method := "AttendanceHandler.Create"
	h.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")

	auth := middleware.GetAuth(ctx)
	request := new(model.CreateAttendanceRequest)
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

	requestCtx := context.WithValue(ctx.UserContext(), "request_id", logger.GetRequestID(ctx))
	err := h.UseCase.Create(requestCtx, request, auth)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: err.Error(),
		})
	}

	h.Log.WithContext(ctx).WithField("method", method).Trace("[END]")

	return ctx.JSON(model.WebResponse[*model.CreateAttendanceResponse]{
		Ok: true,
	})
}
