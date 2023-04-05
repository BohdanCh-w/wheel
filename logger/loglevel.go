package logger

type LogLevel uint8

const (
	Debug = iota
	Info
	Warn
	Error
	Fatal
)
