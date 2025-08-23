package gormlogger

import (
	"context"
	"fmt"
	"time"

	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"gorm.io/gorm/logger"
)

// Adapter to wrap your own Logger into GORM's logger.Interface
type GormLogger struct {
	log      interfaces.Logger
	logLevel logger.LogLevel
}

func New(log interfaces.Logger, level logger.LogLevel) *GormLogger {
	return &GormLogger{
		log:      log,
		logLevel: level,
	}
}

func (gl *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *gl
	newLogger.logLevel = level
	return &newLogger
}

func (gl *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if gl.logLevel >= logger.Info {
		gl.log.Info(fmt.Sprintf(msg, data...))
	}
}

func (gl *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if gl.logLevel >= logger.Warn {
		gl.log.Warn(fmt.Sprintf(msg, data...))
	}
}

func (gl *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if gl.logLevel >= logger.Error {
		gl.log.Error(fmt.Sprintf(msg, data...))
	}
}

func (gl *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if gl.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && gl.logLevel >= logger.Error:
		gl.log.Error("SQL error", "err", err, "sql", sql, "rows", rows, "elapsed", elapsed)
	case gl.logLevel >= logger.Info:
		gl.log.Debug("SQL trace", "sql", sql, "rows", rows, "elapsed", elapsed)
	}
}
