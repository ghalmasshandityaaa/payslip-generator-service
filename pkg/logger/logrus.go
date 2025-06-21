package logger

import (
	"context"
	"os"
	"payslip-generator-service/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func NewLogger(conf *config.Config) *logrus.Logger {
	log := logrus.New()

	log.SetOutput(os.Stdout)
	log.SetReportCaller(false)
	log.SetLevel(logrus.Level(conf.Logger.Level))
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		PrettyPrint:     conf.Logger.Pretty,
	})

	return log
}

func GetRequestID(ctx *fiber.Ctx) string {
	if requestID := ctx.Locals("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

func GetRequestIDFromContext(ctx context.Context) string {
	if requestID := ctx.Value("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

func WithRequestID(logger *logrus.Logger, ctx *fiber.Ctx) *logrus.Entry {
	requestID := GetRequestID(ctx)
	if requestID != "" {
		return logger.WithField("request_id", requestID)
	}
	return logrus.NewEntry(logger)
}

func WithRequestIDFromContext(logger *logrus.Logger, ctx context.Context) *logrus.Entry {
	requestID := GetRequestIDFromContext(ctx)
	if requestID != "" {
		return logger.WithField("request_id", requestID)
	}
	return logrus.NewEntry(logger)
}
