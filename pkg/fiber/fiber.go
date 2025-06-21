package fiber

import (
	"errors"
	"fmt"
	"payslip-generator-service/config"
	"payslip-generator-service/pkg/logger"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func NewFiber(config *config.Config, log *logrus.Logger) *fiber.App {
	// Create context logger for error handler
	contextLogger := logger.NewContextLogger(log)

	var app = fiber.New(fiber.Config{
		AppName:               fmt.Sprintf("%s - %s", config.App.Name, config.App.Version),
		ServerHeader:          "payslip-generator-service",
		ReadTimeout:           time.Duration(config.App.ReadTimeout) * time.Second,
		WriteTimeout:          time.Duration(config.App.WriteTimeout) * time.Second,
		DisableStartupMessage: true,
		ErrorHandler:          errorHandler(contextLogger),
		Prefork:               config.App.Prefork,
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
	})

	return app
}

func errorHandler(log *logger.ContextLogger) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		log.WithContext(ctx).Errorf("request error occurred: %s", err)

		return ctx.Status(code).JSON(fiber.Map{
			"ok":     false,
			"errors": "internal/server-error",
		})
	}
}
