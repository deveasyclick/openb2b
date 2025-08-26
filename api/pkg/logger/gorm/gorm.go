package gormlogger

import (
	"context"
	"fmt"
	"strings"
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
	elapsed := time.Since(begin)
	sql, rows := fc()

	// Always log errors
	if err != nil && gl.logLevel >= logger.Error {
		gl.log.Error("SQL error", "err", err, "sql", sql, "rows", rows, "elapsed", elapsed)
		return
	}

	// Log slow queries (threshold configurable)
	slowThreshold := 200 * time.Millisecond
	if elapsed > slowThreshold && gl.logLevel >= logger.Warn {
		gl.log.Warn("Slow SQL query", "sql", sql, "rows", rows, "elapsed", elapsed)
		return
	}

	// Skip GORM metadata / schema inspection queries
	skipPatterns := []string{
		"SELECT c.column_name", // table constraints
		"SELECT a.attname",     // column types
		"SELECT description FROM pg_catalog.pg_description",
		"SELECT count(*) FROM INFORMATION_SCHEMA.table_constraints",
		"SELECT count(*) FROM pg_indexes",
		"SELECT count(*) FROM information_schema.tables",
		"SELECT CURRENT_DATABASE()",
		"SELECT constraint_name FROM information_schema.table_constraints",
	}
	for _, p := range skipPatterns {
		if strings.HasPrefix(sql, p) {
			return
		}
	}

	// Log all other queries in debug mode
	if gl.logLevel >= logger.Info {
		gl.log.Debug("SQL query", "sql", sql, "rows", rows, "elapsed", elapsed)
	}
}
