package global

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

var DB *gorm.DB

const (
	poolMaxIdleConns = 8
	poolMaxOpenConns = 32
	poolMaxLifetime  = time.Hour
)

func SetupMySQL() {
	logLevel := logger.Warn
	if CFG.Debug {
		logLevel = logger.Info
	}

	var err error
	DB, err = gorm.Open(
		mysql.New(mysql.Config{DSN: CFG.MysqlDSN}),
		&gorm.Config{Logger: &gormLogger{LOG.Handler(), logLevel}})
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

type gormLogger struct {
	handler slog.Handler
	level   logger.LogLevel
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	ll := *l
	ll.level = level
	return &ll
}

func (l *gormLogger) Info(ctx context.Context, format string, args ...any) {
	if l.level >= logger.Info {
		l.logf(ctx, slog.LevelInfo, format, args...)
	}
}

func (l *gormLogger) Warn(ctx context.Context, format string, args ...any) {
	if l.level >= logger.Warn {
		l.logf(ctx, slog.LevelWarn, format, args...)
	}
}

func (l *gormLogger) Error(ctx context.Context, format string, args ...any) {
	if l.level >= logger.Error {
		l.logf(ctx, slog.LevelError, format, args...)
	}
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	sql, rows := fc()
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.level >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound)):
		l.logAttrs(ctx, slog.LevelError, err.Error(),
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.Duration("duration", elapsed),
		)
	case elapsed > 100*time.Millisecond && l.level >= logger.Warn:
		l.logAttrs(ctx, slog.LevelWarn, "slow query",
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.Duration("duration", elapsed),
		)
	case l.level >= logger.Info:
		l.logAttrs(ctx, slog.LevelDebug, "",
			slog.String("sql", sql),
			slog.Int64("rows", rows),
			slog.Duration("duration", elapsed),
		)
	}
}

func (l *gormLogger) logf(ctx context.Context, level slog.Level, format string, args ...any) error {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	pc := pcs[0]
	r := slog.NewRecord(time.Now(), level, fmt.Sprintf(format, args...), pc)
	return l.handler.Handle(ctx, r)
}

func (l *gormLogger) logAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) error {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	pc := pcs[0]
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.AddAttrs(attrs...)
	return l.handler.Handle(ctx, r)
}
