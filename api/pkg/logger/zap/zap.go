package zap

import (
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"go.uber.org/zap"
)

type ZapLogger struct {
	sugar *zap.SugaredLogger
}

func NewZapLogger() interfaces.Logger {
	l, _ := zap.NewProduction()
	return &ZapLogger{sugar: l.Sugar()}
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

func (l *ZapLogger) WithValues(keysAndValues ...interface{}) interfaces.Logger {
	return &ZapLogger{sugar: l.sugar.With(keysAndValues...)}
}
