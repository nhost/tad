package runbook

type RuntimeError struct {
	message string
}

func NewRuntimeError(message string) *RuntimeError {
	return &RuntimeError{
		message: message,
	}
}

func (e *RuntimeError) Error() string {
	return e.message
}
