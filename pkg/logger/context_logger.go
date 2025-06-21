package logger

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ContextLogger struct {
	logger *logrus.Logger
}

func NewContextLogger(logger *logrus.Logger) *ContextLogger {
	return &ContextLogger{
		logger: logger,
	}
}

func (cl *ContextLogger) WithContext(ctx interface{}) *logrus.Entry {
	switch c := ctx.(type) {
	case *fiber.Ctx:
		return WithRequestID(cl.logger, c)
	case context.Context:
		return WithRequestIDFromContext(cl.logger, c)
	default:
		return cl.logger.WithField("error", "invalid context type")
	}
}
