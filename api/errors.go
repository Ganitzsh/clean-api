package api

import "errors"

type APIError struct {
	DataError bool   `json:"-"`
	Err       error  `json:"-"`
	Message   string `json:"error"`
}

func (e APIError) Error() string {
	return e.Err.Error()
}

func NewAPIError(dataError bool, err error) *APIError {
	return &APIError{
		DataError: dataError,
		Err:       err,
	}
}

var (
	ErrSomethingWentWrong     = errors.New("Something went wrong")
	ErrNotImplemented         = errors.New("Not implemented")
	ErrNotFound               = errors.New("Record not found")
	ErrNilValue               = errors.New("Cannot use nil value")
	ErrInvalidInput           = NewAPIError(true, errors.New("Invalid input"))
	ErrUnknownFilterType      = errors.New("Unknown filter type")
	ErrUnsupportedFilterType  = errors.New("Unsupported filter type")
	ErrUnsupportedFilterValue = errors.New("Unsupported filter value")
)
