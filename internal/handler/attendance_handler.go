package handler

import (
	"payslip-generator-service/internal/middleware"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/usecase"
	"payslip-generator-service/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AttendanceHandler struct {
	Log       *logrus.Logger
	UseCase   *usecase.AttendanceUseCase
	Validator *validator.Validator
}

func NewAttendanceHandler(
	useCase *usecase.AttendanceUseCase,
	logger *logrus.Logger,
	validator *validator.Validator,
) *AttendanceHandler {
	return &AttendanceHandler{
		Log:       logger,
		UseCase:   useCase,
		Validator: validator,
	}
}

func (h *AttendanceHandler) Create(ctx *fiber.Ctx) error {
	method := "AttendanceHandler.Create"
	h.Log.Trace("[BEGIN] - ", method)

	auth := middleware.GetAuth(ctx)
	request := new(model.CreateAttendanceRequest)
	if err := ctx.BodyParser(request); err != nil {
		h.Log.Error("failed parse body: ", err.Error())
		return fiber.ErrBadRequest
	}

	errValidation := h.Validator.ValidateStruct(request)
	if errValidation != nil {
		return ctx.JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: errValidation,
		})
	}

	err := h.UseCase.Create(ctx.UserContext(), request, auth)
	if err != nil {
		return ctx.JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: err.Error(),
		})
	}

	h.Log.Trace("[END] - ", method)

	return ctx.JSON(model.WebResponse[*model.CreateAttendanceResponse]{
		Ok: true,
	})
}
