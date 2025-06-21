package gorm

import (
	"context"
	"database/sql"
	"fmt"
	"payslip-generator-service/config"
	"payslip-generator-service/pkg/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/sirupsen/logrus"
)

type GormDB struct {
	db    *gorm.DB
	sqlDB *sql.DB
}

// CustomGormLogger implements gorm.Logger interface with JSON formatting using logrus
type CustomGormLogger struct {
	level  gormlogger.LogLevel
	logger *logrus.Logger
}

func NewCustomGormLogger(log *logrus.Logger) *CustomGormLogger {
	return &CustomGormLogger{
		level:  gormlogger.Info,
		logger: log,
	}
}

func (l *CustomGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

func (l *CustomGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Info {
		entry := logger.WithRequestIDFromContext(l.logger, ctx)
		entry.WithFields(logrus.Fields{
			"component": "gorm",
			"data":      data,
		}).Info(msg)
	}
}

func (l *CustomGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Warn {
		entry := logger.WithRequestIDFromContext(l.logger, ctx)
		entry.WithFields(logrus.Fields{
			"component": "gorm",
			"data":      data,
		}).Warn(msg)
	}
}

func (l *CustomGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Error {
		entry := logger.WithRequestIDFromContext(l.logger, ctx)
		entry.WithFields(logrus.Fields{
			"component": "gorm",
			"data":      data,
		}).Error(msg)
	}
}

func (l *CustomGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	entry := logger.WithRequestIDFromContext(l.logger, ctx)
	fields := logrus.Fields{
		"component":     "gorm",
		"elapsed":       elapsed.String(),
		"rows_affected": rows,
		"sql":           sql,
	}

	if err != nil {
		fields["error"] = err.Error()
		entry.WithFields(fields).Error("gorm query error")
	} else if l.level >= gormlogger.Info {
		entry.WithFields(fields).Info("gorm query executed")
	}
}

func NewGormDB(conf *config.Config, log *logrus.Logger) *GormDB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.Dbname,
		conf.Postgres.SSLMode,
	)

	// Create custom logger for GORM using logrus
	gormLogger := NewCustomGormLogger(log)

	gormConf := &gorm.Config{
		SkipDefaultTransaction: true,
		DryRun:                 conf.Postgres.DryRun,
		PrepareStmt:            true,
		Logger:                 gormLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "public.",
		},
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConf)

	if err != nil {
		panic(fmt.Errorf("failed open database connection: %v", err))
	}

	sqldb, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("connection refused error: %v", err))
	}

	sqldb.SetMaxIdleConns(conf.Postgres.MaxIdleCons)
	sqldb.SetMaxOpenConns(conf.Postgres.MaxOpenCons)
	sqldb.SetConnMaxIdleTime(time.Duration(conf.Postgres.ConnMaxIdleTime) * time.Minute)
	sqldb.SetConnMaxLifetime(time.Duration(conf.Postgres.ConnMaxLifetime) * time.Minute)

	if err := sqldb.Ping(); err != nil {
		panic(fmt.Errorf("ping database got failed: %v", err))
	}

	return &GormDB{db, sqldb}
}

func (g *GormDB) SqlDB() *sql.DB {
	return g.sqlDB
}

func (g *GormDB) DB() *gorm.DB {
	return g.db
}

// Close current db connection. If database connection is not an io.Closer, returns an error.
func (g *GormDB) Close() {
	err := g.sqlDB.Close()

	if err != nil {
		panic(fmt.Errorf("failed close database connection: %v", err))
	}
}
