package logger

type Logger interface {
	Info(msg string, field ...Field)
	Debug(msg string, field ...Field)
	Warn(msg string, field ...Field)
	Error(msg string, field ...Field)
}
