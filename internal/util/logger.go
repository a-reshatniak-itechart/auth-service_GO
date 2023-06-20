package util

import "go.uber.org/zap"

type CustomLogger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type ZapCustomLogger struct {
	Logger zap.Logger
}

func (cl ZapCustomLogger) Info(msg string) {
	cl.Logger.Info(msg)
}

func (cl ZapCustomLogger) Warn(msg string) {
	cl.Logger.Warn(msg)
}

func (cl ZapCustomLogger) Error(msg string) {
	cl.Logger.Error(msg)
}
