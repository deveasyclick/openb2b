package zap

import (
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	sugar *zap.SugaredLogger
}

func NewZapLogger(env string) interfaces.Logger {
	var base *zap.Logger
	var err error

	if env == "production" {
		// JSON logs for production
		cfg := zap.NewProductionConfig()
		cfg.DisableStacktrace = true // only Fatal/Panic get stacktraces
		base, err = cfg.Build()
	} else {
		// Colorful, human-friendly logs
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		base, err = cfg.Build()
	}

	if err != nil {
		panic(err)
	}

	return &ZapLogger{sugar: base.Sugar()}
}

func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}

func (l *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.sugar.Debugw(msg, keysAndValues...)
}

func (l *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.sugar.Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.sugar.Fatalw(msg, keysAndValues...)
}

func (l *ZapLogger) WithValues(keysAndValues ...interface{}) interfaces.Logger {
	return &ZapLogger{sugar: l.sugar.With(keysAndValues...)}
}
