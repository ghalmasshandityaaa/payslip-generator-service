package gorm

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"payslip-generator-service/config"
	"time"
)

type GormDB struct {
	db    *gorm.DB
	sqlDB *sql.DB
}

func NewGormDB(conf *config.Config) *GormDB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.Dbname,
		conf.Postgres.SSLMode,
	)

	gormConf := &gorm.Config{
		SkipDefaultTransaction: true,
		DryRun:                 conf.Postgres.DryRun,
		PrepareStmt:            true,
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
