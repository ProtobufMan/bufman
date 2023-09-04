package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	once   sync.Once
	logger *logrus.Logger
)

func initLogger() {
	if logger == nil {
		once.Do(func() {
			if logger == nil {
				logger = logrus.New()
			}
		})
	}
}

func SetLevel(mode string) error {
	initLogger()

	switch mode {
	case gin.DebugMode:
		logger.SetLevel(logrus.DebugLevel)
	case gin.TestMode:
		logger.SetLevel(logrus.InfoLevel)
	case gin.ReleaseMode:
		logger.SetLevel(logrus.ErrorLevel)
	default:
		return fmt.Errorf("not support mode for %s", mode)
	}

	return nil
}

type LogLevel uint32

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelPanic
	LogLevelFatal
)

func log(level LogLevel, args ...interface{}) {
	initLogger()

	switch level {
	case LogLevelDebug:
		logger.Debug(args...)
	case LogLevelInfo:
		logger.Info(args...)
	case LogLevelWarn:
		logger.Warn(args...)
	case LogLevelError:
		logger.Error(args...)
	case LogLevelPanic:
		logger.Panic(args...)
	case LogLevelFatal:
		logger.Fatal(args...)
	}
}

func logFormat(level LogLevel, format string, args ...interface{}) {
	initLogger()

	switch level {
	case LogLevelDebug:
		logger.Debugf(format, args...)
	case LogLevelInfo:
		logger.Infof(format, args...)
	case LogLevelWarn:
		logger.Warnf(format, args...)
	case LogLevelError:
		logger.Errorf(format, args...)
	case LogLevelPanic:
		logger.Panicf(format, args...)
	case LogLevelFatal:
		logger.Fatalf(format, args...)
	}
}

func Debug(args ...interface{}) {
	log(LogLevelDebug, args)
}

func Debugf(format string, args ...interface{}) {
	logFormat(LogLevelDebug, format, args...)
}

func Info(args ...interface{}) {
	log(LogLevelInfo, args)
}

func Infof(format string, args ...interface{}) {
	logFormat(LogLevelInfo, format, args...)
}

func Warn(args ...interface{}) {
	log(LogLevelWarn, args)
}

func Warnf(format string, args ...interface{}) {
	logFormat(LogLevelWarn, format, args...)
}

func Error(args ...interface{}) {
	log(LogLevelError, args)
}

func Errorf(format string, args ...interface{}) {
	logFormat(LogLevelError, format, args...)
}

func Panic(args ...interface{}) {
	log(LogLevelPanic, args)
}

func Panicf(format string, args ...interface{}) {
	logFormat(LogLevelPanic, format, args)
}

func Fatal(args ...interface{}) {
	log(LogLevelFatal, args)
}

func Fatalf(format string, args ...interface{}) {
	logFormat(LogLevelFatal, format, args)
}
