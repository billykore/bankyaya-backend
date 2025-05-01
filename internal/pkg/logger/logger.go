package logger

import (
	"fmt"
	"os"
	"sync"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

// Logger provides logging functions.
type Logger struct {
	zapLogger *zap.Logger
	domain    string
	usecase   string
}

var (
	_once sync.Once
	_log  *Logger
)

// New return the singleton instance of Logger.
func New() *Logger {
	_once.Do(func() {
		zapLogger := newLogger()
		_log = &Logger{zapLogger: zapLogger}
	})
	return _log
}

func newLogger() *zap.Logger {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.InfoLevel)
	logger := zap.New(core, zap.AddStacktrace(zap.PanicLevel))
	return logger
}

func (l *Logger) DomainUsecase(domain, usecase string) *Logger {
	l.domain = domain
	l.usecase = usecase
	return l
}

// Error log.
func (l *Logger) Error(err error) {
	l.zapLogger.Error(
		err.Error(),
		zap.String("domain", l.domain),
		zap.String("usecase", l.usecase),
	)
}

// Errorf is formated Error log.
func (l *Logger) Errorf(format string, a ...any) {
	l.zapLogger.Error(
		fmt.Sprintf(format, a...),
		zap.String("domain", l.usecase),
		zap.String("usecase", l.usecase),
	)
}

// Fatal logs an error then shutdown the domain.
func (l *Logger) Fatal(err error) {
	l.zapLogger.Fatal(
		err.Error(),
		zap.String("usecase", l.usecase),
	)
}

// Fatalf is formated Fatal log.
func (l *Logger) Fatalf(format string, a ...any) {
	l.zapLogger.Fatal(
		fmt.Sprintf(format, a...),
		zap.String("usecase", l.usecase),
	)
}

// Info log.
func (l *Logger) Info(v any) {
	l.zapLogger.Info(
		fmt.Sprintf("%v", v),
		zap.String("domain", l.domain),
		zap.String("usecase", l.usecase),
	)
}

// Infof is formated Info log.
func (l *Logger) Infof(format string, a ...any) {
	l.zapLogger.Info(
		fmt.Sprintf(format, a...),
		zap.String("domain", l.domain),
		zap.String("usecase", l.usecase),
	)
}
