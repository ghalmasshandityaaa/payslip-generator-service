package handler

import (
	"payslip-generator-service/internal/middleware"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/usecase"
	"payslip-generator-service/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ReimbursementHandler struct {
	Log       *logrus.Logger
	UseCase   *usecase.ReimbursementUseCase
	Validator *validator.Validator
}

func NewReimbursementHandler(
	useCase *usecase.ReimbursementUseCase,
	logger *logrus.Logger,
	validator *validator.Validator,
) *ReimbursementHandler {
	return &ReimbursementHandler{
		Log:       logger,
		UseCase:   useCase,
		Validator: validator,
	}
}

// Create creates a new reimbursement record for the authenticated employee
// @Summary Create reimbursement record
// @Description Create a new reimbursement record with amount and description for the authenticated employee
// @Tags Reimbursement
// @Accept json
// @Produce json
// @Security bearer
// @Param request body model.CreateReimbursementRequest true "Reimbursement details"
// @Router /reimbursement [post]
func (h *ReimbursementHandler) Create(ctx *fiber.Ctx) error {
	method := "ReimbursementHandler.Create"
	h.Log.Trace("[BEGIN] - ", method)

	auth := middleware.GetAuth(ctx)
	request := new(model.CreateReimbursementRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}

	errValidation := h.Validator.ValidateStruct(request)
	if errValidation != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: errValidation,
		})
	}

	err := h.UseCase.Create(ctx.UserContext(), request, auth)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: err.Error(),
		})
	}

	h.Log.Trace("[END] - ", method)
	return ctx.JSON(model.WebResponse[*model.CreateReimbursementResponse]{
		Ok: true,
	})
}
