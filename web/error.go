package web

type WebError struct { // nolint: revive
	Code int   `json:"status"`
	Err  error `json:"error"`
}

func NewError(code int, err error) error {
	return &WebError{
		Code: code,
		Err:  err,
	}
}

func (e *WebError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return ""
}

func (e *WebError) Status() int {
	return e.Code
}

func (e *WebError) Unwrap() error {
	return e.Err
}
