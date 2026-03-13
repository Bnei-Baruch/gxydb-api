package errs

type WithMessage struct {
	Msg string
	Err error
}

func (e *WithMessage) Error() string {
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

func (e *WithMessage) Cause() error {
	return e.Err
}
