package logger

import (
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/deveasyclick/openb2b/pkg/logger/zap"
)

func New() interfaces.Logger {
	return zap.NewZapLogger()
}
