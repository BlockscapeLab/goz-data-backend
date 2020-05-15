package types

// ErrorNotFound is an error type for NOT_FOUND messages and will cause response status code to be 404
type ErrorNotFound struct {
	message string
}

func (e ErrorNotFound) Error() string {
	return e.message
}

// NewErrorNotFound constructor
func NewErrorNotFound(message string) ErrorNotFound {
	return ErrorNotFound{message}
}

// GetStatusCode will return status code that matches error if error from this rest package is used. Will return 500 on unknown errors.
func GetStatusCode(err error) int {
	switch err.(type) {
	case ErrorNotFound:
		return 404
	default:
		return 500
	}
}
