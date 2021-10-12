package web

type ErrorResponse struct {
	Error string `json:"error"`
}

type Error struct {
	Err    error
	Status int
}

func NewRequestError(err error, status int) error {
	return &Error{Err: err, Status: status}
}

func (err *Error) Error() string {
	return err.Err.Error()
}
