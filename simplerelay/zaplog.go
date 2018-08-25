package main

import (
	"log"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	zaplog *zap.SugaredLogger
)

type ZapLoger struct {
}

func (*ZapLoger) Trace(format string, v ...interface{}) {
	if zaplog != nil {
		zaplog.Debugf(format, v...)
	} else {
		log.Printf(format, v...)
	}
}
func (*ZapLoger) Debug(format string, v ...interface{}) {
	if zaplog != nil {
		zaplog.Debugf(format, v...)
	} else {
		log.Printf(format, v...)
	}
}
func (*ZapLoger) Info(format string, v ...interface{}) {
	if zaplog != nil {
		zaplog.Infof(format, v...)
	} else {
		log.Printf(format, v...)
	}
}
func (*ZapLoger) Warn(format string, v ...interface{}) {
	if zaplog != nil {
		zaplog.Warnf(format, v...)
	} else {
		log.Printf(format, v...)
	}
}
func (*ZapLoger) Error(format string, v ...interface{}) {
	if zaplog != nil {
		zaplog.Errorf(format, v...)
	} else {
		log.Printf(format, v...)
	}
}
func (*ZapLoger) Faltal(format string, v ...interface{}) {
	if zaplog != nil {
		zaplog.Fatalf(format, v...)
	} else {
		log.Printf(format, v...)
	}
}
