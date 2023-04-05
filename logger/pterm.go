package logger

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pterm/pterm"
)

func ptermLevel(level LogLevel) pterm.LogLevel {
	switch level {
	case Debug:
		return pterm.LogLevelDebug
	case Info:
		return pterm.LogLevelInfo
	case Warn:
		return pterm.LogLevelWarn
	case Error:
		return pterm.LogLevelError
	case Fatal:
		return pterm.LogLevelFatal
	default:
		return pterm.LogLevelInfo
	}
}

var _ Logger = (*PtermLogger)(nil)

func NewPtermLogger(level LogLevel) *PtermLogger {
	return &PtermLogger{
		log: pterm.DefaultLogger.WithLevel(ptermLevel(level)),
	}
}

type PtermLogger struct {
	log  *pterm.Logger
	tID  uuid.UUID
	args map[string]any
}

func (ptl *PtermLogger) WithLevel(level LogLevel) Logger {
	ptl.log.Level = ptermLevel(level)

	return ptl
}

func (ptl *PtermLogger) WithTransaction(id uuid.UUID) Logger {
	ptl.tID = id

	return ptl
}

func (ptl *PtermLogger) With(key string, value any) Logger {
	ptl.args[key] = value

	return ptl
}

func (ptl *PtermLogger) Debugf(msg string, args ...any) {
	ptl.log.Debug(fmt.Sprintf(msg, args...), ptl.argumets())
}

func (ptl *PtermLogger) Infof(msg string, args ...any) {
	ptl.log.Info(fmt.Sprintf(msg, args...), ptl.argumets())
}

func (ptl *PtermLogger) Warnf(msg string, args ...any) {
	ptl.log.Warn(fmt.Sprintf(msg, args...), ptl.argumets())
}

func (ptl *PtermLogger) Errorf(msg string, args ...any) {
	ptl.log.Error(fmt.Sprintf(msg, args...), ptl.argumets())
}

func (ptl *PtermLogger) Fatalf(msg string, args ...any) {
	ptl.log.Fatal(fmt.Sprintf(msg, args...), ptl.argumets())
}

func (ptl *PtermLogger) argumets() []pterm.LoggerArgument {
	args := make([]pterm.LoggerArgument, 0, len(ptl.args)+1)

	if ptl.tID != uuid.Nil {
		args = append(args, ptl.log.Args(TransationKey, ptl.tID)...)
	}

	return append(args, ptl.log.ArgsFromMap(ptl.args)...)
}
