package handler

import (
	"context"
	"math"
	"payslip-generator-service/internal/entity"
	"payslip-generator-service/internal/middleware"
	"payslip-generator-service/internal/model"
	"payslip-generator-service/internal/usecase"
	"payslip-generator-service/internal/vm"
	"payslip-generator-service/pkg/logger"
	"payslip-generator-service/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type PayrollHandler struct {
	Log       *logger.ContextLogger
	UseCase   *usecase.PayrollUseCase
	Validator *validator.Validator
}

func NewPayrollHandler(
	useCase *usecase.PayrollUseCase,
	log *logger.ContextLogger,
	validator *validator.Validator,
) *PayrollHandler {
	return &PayrollHandler{
		Log:       log,
		UseCase:   useCase,
		Validator: validator,
	}
}

// ListPeriod retrieves a paginated list of payroll periods
// @Summary List payroll periods
// @Description Get a paginated list of all payroll periods in the system
// @Tags Payroll
// @Accept json
// @Produce json
// @Security bearer
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param size query int false "Page size (default: 10)" minimum(1)
// @Router /payroll/period [get]
func (h *PayrollHandler) ListPeriod(ctx *fiber.Ctx) error {
	method := "PayrollHandler.ListPeriod"
	h.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")

	request := &model.ListPayrollPeriodRequest{
		Page:     ctx.QueryInt("page", 1),
		PageSize: ctx.QueryInt("size", 10),
	}

	// Create context with request_id
	requestCtx := context.WithValue(ctx.UserContext(), "request_id", logger.GetRequestID(ctx))
	data, total, err := h.UseCase.ListPeriod(requestCtx, request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
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

	h.Log.WithContext(ctx).WithField("method", method).Trace("[END]")
	return ctx.JSON(model.WebResponseWithData[[]entity.PayrollPeriod]{
		Ok:     true,
		Data:   data,
		Paging: paging,
	})
}

// CreatePeriod creates a new payroll period
// @Summary Create payroll period
// @Description Create a new payroll period with start and end dates (Admin only)
// @Tags Payroll
// @Accept json
// @Produce json
// @Security bearer
// @Param request body model.CreatePayrollPeriodRequest true "Payroll period details"
// @Router /payroll/period [post]
func (h *PayrollHandler) CreatePeriod(ctx *fiber.Ctx) error {
	method := "PayrollHandler.Create"
	h.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")

	auth := middleware.GetAuth(ctx)
	request := new(model.CreatePayrollPeriodRequest)
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

	// Create context with request_id
	requestCtx := context.WithValue(ctx.UserContext(), "request_id", logger.GetRequestID(ctx))
	err := h.UseCase.CreatePeriod(requestCtx, request, auth)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: err.Error(),
		})
	}

	h.Log.WithContext(ctx).WithField("method", method).Trace("[END]")
	return ctx.JSON(model.WebResponse[*model.CreatePayrollPeriodResponse]{
		Ok: true,
	})
}

// ProcessPayroll processes payroll for a specific period
// @Summary Process payroll
// @Description Process payroll calculations for all employees in a specific period (Admin only)
// @Tags Payroll
// @Accept json
// @Produce json
// @Security bearer
// @Param request body model.ProcessPayrollRequest true "Payroll processing details"
// @Router /payroll/process [post]
func (h *PayrollHandler) ProcessPayroll(ctx *fiber.Ctx) error {
	method := "PayrollHandler.ProcessPayroll"
	h.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")

	auth := middleware.GetAuth(ctx)
	request := new(model.ProcessPayrollRequest)
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

	// Create context with request_id
	requestCtx := context.WithValue(ctx.UserContext(), "request_id", logger.GetRequestID(ctx))
	err := h.UseCase.ProcessPayroll(requestCtx, request, auth)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: err.Error(),
		})
	}

	h.Log.WithContext(ctx).WithField("method", method).Trace("[END]")
	return ctx.JSON(model.WebResponse[*model.CreatePayrollPeriodResponse]{
		Ok: true,
	})
}

// GetPayslip retrieves payslip for the authenticated employee
// @Summary Get payslip
// @Description Get payslip details for the authenticated employee in a specific period (Employee only)
// @Tags Payroll
// @Accept json
// @Produce json
// @Security bearer
// @Param period_id query string true "Payroll period ID" example("01HXYZ123456789ABCDEFGHIJK")
// @Router /payroll/payslip [get]
func (h *PayrollHandler) GetPayslip(ctx *fiber.Ctx) error {
	method := "PayrollHandler.GetPayslip"
	h.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")

	auth := middleware.GetAuth(ctx)

	request := &model.GetPayslipRequest{
		PeriodID: ctx.Query("period_id"),
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
	data, err := h.UseCase.GetPayslip(requestCtx, request, auth)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: err.Error(),
		})
	}

	h.Log.WithContext(ctx).WithField("method", method).Trace("[END]")

	return ctx.JSON(model.WebResponse[*vm.Payslip]{
		Ok:   true,
		Data: data,
	})
}

// GetPayslipReport retrieves payslip report for all employees
// @Summary Get payslip report
// @Description Get comprehensive payslip report for all employees in a specific period (Admin only)
// @Tags Payroll
// @Accept json
// @Produce json
// @Security bearer
// @Param period_id query string true "Payroll period ID" example("01HXYZ123456789ABCDEFGHIJK")
// @Router /payroll/payslip/report [get]
func (h *PayrollHandler) GetPayslipReport(ctx *fiber.Ctx) error {
	method := "PayrollHandler.GetPayslipReport"
	h.Log.WithContext(ctx).WithField("method", method).Trace("[BEGIN]")

	auth := middleware.GetAuth(ctx)

	request := &model.GetPayslipRequest{
		PeriodID: ctx.Query("period_id"),
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
	data, err := h.UseCase.GetPayslipReport(requestCtx, request, auth)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
			Ok:     false,
			Errors: err.Error(),
		})
	}

	h.Log.WithContext(ctx).WithField("method", method).Trace("[END]")

	return ctx.JSON(model.WebResponse[*vm.PayslipReport]{
		Ok:   true,
		Data: data,
	})
}
