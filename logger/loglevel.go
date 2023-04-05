package logger

import (
	"fmt"
	"strings"

	wherr "github.com/bohdanch-w/wheel/errors"
)

type LogLevel uint8

const (
	Debug = iota
	Info
	Warn
	Error
	Fatal

	invalid
	errInvalidLogLevel = wherr.Error("loglevel is invalid")
)

func LogLevelFromString(str string) LogLevel {
	switch strings.ToLower(strings.TrimSpace(str)) {
	case "debug":
		return Debug
	case "info":
		return Info
	case "warn":
		return Warn
	case "error":
		return Error
	case "fatal":
		return Fatal
	case "":
		return Info
	}

	return invalid
}

func (l *LogLevel) UnmarshalText(text []byte) error {
	strLevel := string(text)

	level := LogLevelFromString(strLevel)
	if level == invalid {
		return fmt.Errorf("%w: %q", errInvalidLogLevel, strLevel)
	}

	*l = level

	return nil
}
