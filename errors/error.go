package errors

import "fmt"

type Error string

func (e Error) Error() string {
	return string(e)
}

func Errorf(format, err string, args ...any) error {
	return fmt.Errorf(format, append([]any{Error(err)}, args...)...) // nolint: err113
}
