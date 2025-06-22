package logger

import (
	"context"
	"os"
	"payslip-generator-service/config"
	"strings"
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

// GetIPAddress extracts IP address from Fiber context
func GetIPAddress(ctx *fiber.Ctx) string {
	// Try to get IP from X-Forwarded-For header first (for proxy scenarios)
	if forwardedFor := ctx.Get("X-Forwarded-For"); forwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if commaIndex := strings.Index(forwardedFor, ","); commaIndex != -1 {
			return strings.TrimSpace(forwardedFor[:commaIndex])
		}
		return strings.TrimSpace(forwardedFor)
	}

	// Try X-Real-IP header
	if realIP := ctx.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Fallback to remote IP
	return ctx.IP()
}

// GetIPAddressFromContext extracts IP address from context
func GetIPAddressFromContext(ctx context.Context) string {
	if ip := ctx.Value("ip_address"); ip != nil {
		if ipStr, ok := ip.(string); ok {
			return ipStr
		}
	}
	return ""
}

// GetUserAgent extracts User-Agent from Fiber context
func GetUserAgent(ctx *fiber.Ctx) string {
	return ctx.Get("User-Agent")
}

// GetUserAgentFromContext extracts User-Agent from context
func GetUserAgentFromContext(ctx context.Context) string {
	if ua := ctx.Value("user_agent"); ua != nil {
		if uaStr, ok := ua.(string); ok {
			return uaStr
		}
	}
	return ""
}

func WithRequestID(logger *logrus.Logger, ctx *fiber.Ctx) *logrus.Entry {
	requestID := GetRequestID(ctx)
	ipAddress := GetIPAddress(ctx)
	userAgent := GetUserAgent(ctx)

	entry := logrus.NewEntry(logger)
	if requestID != "" {
		entry = entry.WithField("request_id", requestID)
	}
	if ipAddress != "" {
		entry = entry.WithField("ip_address", ipAddress)
	}
	if userAgent != "" {
		entry = entry.WithField("user_agent", userAgent)
	}
	return entry
}

func WithRequestIDFromContext(logger *logrus.Logger, ctx context.Context) *logrus.Entry {
	requestID := GetRequestIDFromContext(ctx)
	ipAddress := GetIPAddressFromContext(ctx)
	userAgent := GetUserAgentFromContext(ctx)

	entry := logrus.NewEntry(logger)
	if requestID != "" {
		entry = entry.WithField("request_id", requestID)
	}
	if ipAddress != "" {
		entry = entry.WithField("ip_address", ipAddress)
	}
	if userAgent != "" {
		entry = entry.WithField("user_agent", userAgent)
	}
	return entry
}
