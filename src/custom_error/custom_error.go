package custom_error

type AppError struct {
	Message       string
	HttpErrorCode int
	Err           error
}

func (e *AppError) Error() string {
	return e.Message
}
