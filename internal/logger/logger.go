package logger

import (
	"dndroller/internal/config"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = "debug"
	ErrorLevel = "error"
	InfoLevel  = "info"
	WarnLevel  = "warn"
)

type Zap struct {
	*zap.SugaredLogger
}

func NewZapLogger(cfg *config.Config) (*Zap, error) {
	config := zap.NewProductionConfig()
	config.Level = zapLevel(cfg.LogLevel)
	config.Encoding = "json"
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.DisableCaller = true
	baseLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Zap{baseLogger.Sugar()}, nil
}

func (that Zap) Println(elems ...interface{}) {
	that.Debug(elems...)
	that.Debug("\n")
}

func (that Zap) Printf(format string, elems ...interface{}) {
	that.Debugf(format, elems...)
}

func (that Zap) Print(params ...interface{}) {
	that.Debug(params...)
}

func (that Zap) LogAndWrapError(err error, message string) error {
	e := fmt.Errorf("%s:%v", message, err)
	that.Error(e.Error())
	return e
}

func zapLevel(level string) (l zap.AtomicLevel) {
	switch level {
	case DebugLevel:
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case ErrorLevel:
		return zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case InfoLevel:
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case WarnLevel:
		return zap.NewAtomicLevelAt(zapcore.WarnLevel)
	default:
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
}
