package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"payslip-generator-service/config"
	"time"
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
