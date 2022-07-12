package values

import (
	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

type LogLevel struct {
	zapcore.Level
}

var _ pflag.Value = &LogLevel{}

func (l LogLevel) Type() string {
	return "LogLevel"
}
