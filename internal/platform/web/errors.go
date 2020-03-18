package web

type ErrorResponse struct {
	Error string `json:"error"`
}

// Error is sued to add web info to request error.
type Error struct {
	Err    error
	Status int
}

func NewRequestError(err error, status int) error {
	return &Error{Err: err, Status: status}
}

func (e *Error) Error() string {
	return e.Err.Error()
}
