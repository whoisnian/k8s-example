package global

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"moul.io/zapgorm2"
)

var DB *gorm.DB

const (
	poolMaxIdleConns = 8
	poolMaxOpenConns = 32
	poolMaxLifetime  = time.Hour
)

func SetupMySQL() {
	gormLogger := zapgorm2.Logger{
		ZapLogger:                 LOG,
		LogLevel:                  logger.Warn,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
		Context:                   nil,
	}
	if CFG.Debug {
		gormLogger.LogLevel = logger.Info
	}

	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN: CFG.MysqlDSN,
	}), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic(err)
	}

	if err := DB.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(poolMaxIdleConns)
	sqlDB.SetMaxOpenConns(poolMaxOpenConns)
	sqlDB.SetConnMaxLifetime(poolMaxLifetime)

	if err = sqlDB.Ping(); err != nil {
		panic(err)
	}
}
