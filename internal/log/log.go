package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//goland:noinspection ALL
var (
	slogger *zap.SugaredLogger

	// format functions
	Debugf, Infof, Warnf, Errorf, Fatalf, Panicf func(string, ...interface{})
	// key values
	Debugw, Infow, Warnw, Errorw, Fatalw, Panicw func(string, ...interface{})
	// With construct new logger with args
	With func(args ...interface{}) *zap.SugaredLogger
)

type logConfig struct {
	config  zap.Config
	options []zap.Option
}

type Option func(*logConfig)

func WithLogLevel(l zapcore.Level) Option {
	return func(c *logConfig) {
		c.config.Level.SetLevel(l)
	}
}

func Init(opts ...Option) {
	c := logConfig{
		config:  zap.NewProductionConfig(),
		options: []zap.Option{},
	}
	for _, opt := range opts {
		opt(&c)
	}

	if l, err := c.config.Build(c.options...); err != nil {
		panic(err)
	} else {
		initLogger(l)
	}
}

func init() {
	Init()
}

func initLogger(l *zap.Logger) {
	slogger = l.Sugar()
	Debugf = slogger.Debugf
	Infof = slogger.Infof
	Warnf = slogger.Warnf
	Errorf = slogger.Errorf
	Panicf = slogger.Panicf

	Debugw = slogger.Debugw
	Infow = slogger.Infow
	Warnw = slogger.Warnw
	Errorw = slogger.Errorw
	Panicw = slogger.Panicw
}
