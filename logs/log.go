package logs

import (
	"log"

	"go.uber.org/zap"
)

type Logger interface {
	Trace(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Faltal(format string, v ...interface{})
}

type DefaultLoger struct {
}

func (*DefaultLoger) Trace(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (*DefaultLoger) Debug(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (*DefaultLoger) Info(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (*DefaultLoger) Warn(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (*DefaultLoger) Error(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (*DefaultLoger) Faltal(format string, v ...interface{}) {
	log.Printf(format, v...)
}

var (
	_log        Logger
	_defaultLog DefaultLoger
	_zapBuilder func(options ...zap.Option) *zap.Logger
	IsZapReady  bool
)

func init() {
	_log = &_defaultLog
	IsZapReady = false
}

func SetLoger(l Logger) {
	_log = l
}

func SetZapLogerBuilder(buildZap func(options ...zap.Option) *zap.Logger) {
	_zapBuilder = buildZap
	IsZapReady = true
}

func GetZapLoger(name string, options ...zap.Option) *zap.Logger {
	return _zapBuilder(options...).Named(name)
}

func Trace(format string, v ...interface{}) {
	if _log != nil {
		_log.Trace(format, v...)
	}
}

func Debug(format string, v ...interface{}) {
	if _log != nil {
		_log.Debug(format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if _log != nil {
		_log.Info(format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if _log != nil {
		_log.Warn(format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if _log != nil {
		_log.Error(format, v...)
	}
}

func Faltal(format string, v ...interface{}) {
	if _log != nil {
		_log.Faltal(format, v...)
	}
}
