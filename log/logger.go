package log

import (
	"fmt"
	"github.com/jslowik/commacloner/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/url"
	"os"
	"path/filepath"
)

var (
	globalLogger *zap.Logger
)

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

// InitWithConfiguration initializes the logger with the set default log level and format.
func InitWithConfiguration(config config.Logger) error {
	var lvl zap.AtomicLevel
	err := lvl.UnmarshalText([]byte(config.Level))
	if err != nil {
		return fmt.Errorf("invalid log level: %v", err)
	}

	cnfg := zap.NewProductionConfig()
	cnfg.Level = lvl
	cnfg.Encoding = config.Format
	cnfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cnfg.OutputPaths = []string{"stderr"}

	if config.Destination == "file" {
		logFile := "./logs/commacloner.log"
		e := createDir("./logs")
		if e != nil {
			return fmt.Errorf("could not create log directory: %v", e)
		}
		ll := lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    1024, //MB
			MaxBackups: 5,
			MaxAge:     90, //days
			Compress:   true,
		}
		zap.RegisterSink("lumberjack", func(*url.URL) (zap.Sink, error) {
			return lumberjackSink{
				Logger: &ll,
			}, nil
		})
		cnfg.OutputPaths = append(cnfg.OutputPaths, fmt.Sprintf("lumberjack:%s", logFile))
	}

	_globalLogger, err := cnfg.Build()
	if err != nil {
		panic(fmt.Sprintf("build zap logger from config error: %v", err))
	}
	zap.ReplaceGlobals(_globalLogger)
	globalLogger = _globalLogger
	return nil
}

// IsDir returns true if the given path is an existing directory.
func createDir(path string) error {
	if pathAbs, err := filepath.Abs(path); err == nil {
		if fileInfo, err := os.Stat(pathAbs); !os.IsNotExist(err) && fileInfo.IsDir() {
			return nil
		} else if e := os.MkdirAll(path, os.ModePerm); e != nil {
			return e
		}
	}
	return nil
}

func NewLogger(name string) *zap.SugaredLogger {
	if globalLogger == nil {
		c := config.Logger{
			Level:       "debug",
			Format:      "console",
			Destination: "console",
		}
		InitWithConfiguration(c)
	}
	return globalLogger.Named(name).Sugar()
}
