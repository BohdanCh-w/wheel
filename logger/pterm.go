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
		log:  pterm.DefaultLogger.WithLevel(ptermLevel(level)),
		args: make(map[string]any),
	}
}

type PtermLogger struct {
	log  *pterm.Logger
	tID  uuid.UUID
	args map[string]any
}

func copyArgs(args map[string]any) map[string]any {
	ret := make(map[string]any, len(args))
	for k, v := range args {
		ret[k] = v
	}

	return ret
}

func (ptl *PtermLogger) WithLevel(level LogLevel) Logger {
	return &PtermLogger{
		log:  ptl.log.WithLevel(ptermLevel(level)),
		tID:  ptl.tID,
		args: copyArgs(ptl.args),
	}
}

func (ptl *PtermLogger) WithTransaction(id uuid.UUID) Logger {
	log := *ptl.log

	return &PtermLogger{
		log:  &log,
		tID:  id,
		args: copyArgs(ptl.args),
	}
}

func (ptl *PtermLogger) WithError(err error) Logger {
	log := *ptl.log
	args := copyArgs(ptl.args)
	args["error"] = err.Error()

	return &PtermLogger{
		log:  &log,
		args: args,
	}
}

func (ptl *PtermLogger) With(key string, value any) Logger {
	log := *ptl.log
	args := copyArgs(ptl.args)
	args[key] = value

	return &PtermLogger{
		log:  &log,
		tID:  ptl.tID,
		args: args,
	}
}

func (ptl *PtermLogger) Debugf(msg string, args ...any) {
	ptl.log.Debug(fmt.Sprintf(msg, args...), ptl.arguments())
}

func (ptl *PtermLogger) Infof(msg string, args ...any) {
	ptl.log.Info(fmt.Sprintf(msg, args...), ptl.arguments())
}

func (ptl *PtermLogger) Warnf(msg string, args ...any) {
	ptl.log.Warn(fmt.Sprintf(msg, args...), ptl.arguments())
}

func (ptl *PtermLogger) Errorf(msg string, args ...any) {
	ptl.log.Error(fmt.Sprintf(msg, args...), ptl.arguments())
}

func (ptl *PtermLogger) Fatalf(msg string, args ...any) {
	ptl.log.Fatal(fmt.Sprintf(msg, args...), ptl.arguments())
}

func (ptl *PtermLogger) arguments() []pterm.LoggerArgument {
	args := make([]pterm.LoggerArgument, 0, len(ptl.args)+1)

	if ptl.tID != uuid.Nil {
		args = append(args, ptl.log.Args(TransactionKey, ptl.tID)...)
	}

	return append(args, ptl.log.ArgsFromMap(ptl.args)...)
}
