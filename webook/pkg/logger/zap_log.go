package logger

import "go.uber.org/zap"

type zap_logger struct {
	l *zap.Logger
}

func (z zap_logger) Info(msg string, field ...Field) {
	z.l.Info(msg, z.toZapField(field)...)
}

func (z zap_logger) Debug(msg string, field ...Field) {
	z.l.Debug(msg, z.toZapField(field)...)
}

func (z zap_logger) Warn(msg string, field ...Field) {
	z.l.Warn(msg, z.toZapField(field)...)
}

func (z zap_logger) Error(msg string, field ...Field) {
	z.l.Error(msg, z.toZapField(field)...)
}

func NewZapLogger(l *zap.Logger) *zap_logger {
	return &zap_logger{l: l}
}

func (z zap_logger) toZapField(field []Field) []zap.Field {
	res := make([]zap.Field, len(field))
	for _, v := range field {
		res = append(res, zap.Any(v.Key, v.Val))
	}
	return res
}
