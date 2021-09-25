package log

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once        sync.Once
	instance    *zap.Logger
	serverError error
)

// InitWithConfiguration initializes the logger with the set default log level and format.
func InitWithConfiguration(level, format string) (*zap.Logger, error) {
	once.Do(func() {
		instance, serverError = newLogger(level, format)
	})
	return instance, serverError
}

func newLogger(level string, format string) (*zap.Logger, error) {
	var lvl zap.AtomicLevel
	err := lvl.UnmarshalText([]byte(level))
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %v", err)
	}

	cnfg := zap.NewProductionConfig()
	cnfg.Level = lvl
	cnfg.Encoding = format
	cnfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cnfg.Build()
	if err != nil {
		panic(err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}(logger)
	return logger, nil
}
