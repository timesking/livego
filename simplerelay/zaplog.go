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
	L *zap.Logger
	S *zap.SugaredLogger
}

func NewRelayLoger() *ZapLoger {
	newLoger := logger.WithOptions(zap.AddCallerSkip(2)).Named("RelayLoger")
	z := &ZapLoger{
		L: newLoger,
		S: newLoger.Sugar(),
	}

	return z
}

func (z *ZapLoger) Trace(format string, v ...interface{}) {
	if z != nil {
		z.S.Debugf(format, v...)
	} else {
		log.Printf(format, v...)
	}
}
func (z *ZapLoger) Debug(format string, v ...interface{}) {
	if z != nil {
		z.S.Debugf(format, v...)
	} else {
		log.Printf(format, v...)
	}
}
func (z *ZapLoger) Info(format string, v ...interface{}) {
	if z != nil {
		z.S.Infof(format, v...)
	} else {
		log.Printf(format, v...)
	}
}
func (z *ZapLoger) Warn(format string, v ...interface{}) {
	if z != nil {
		z.S.Warnf(format, v...)
	} else {
		log.Printf(format, v...)
	}
}
func (z *ZapLoger) Error(format string, v ...interface{}) {
	if z != nil {
		z.S.Errorf(format, v...)
	} else {
		log.Printf(format, v...)
	}
}
func (z *ZapLoger) Faltal(format string, v ...interface{}) {
	if z != nil {
		z.S.Fatalf(format, v...)
	} else {
		log.Printf(format, v...)
	}
}
