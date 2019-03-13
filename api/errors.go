package api

import (
	"errors"
	"net/http"
)

// ErrorCode is a standarized string that identifies issues across the API
type ErrorCode string

// APIError is the content that is returned on error
type APIError struct {
	Message string    `json:"error"`
	AppCode ErrorCode `json:"code,omitempty"`

	DataError  bool  `json:"-"`
	StatusCode int   `json:"-"`
	Err        error `json:"-"`
}

func (e APIError) Error() string {
	ret := e.Message
	if e.Err != nil {
		ret += ": " + e.Err.Error()
	}
	return ret
}

func (e *APIError) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewAPIError(dataError bool, message string) *APIError {
	return &APIError{
		DataError: dataError,
		Message:   message,
	}
}

const (
	ErrorCodeInternalError  ErrorCode = "internal_error"
	ErrorCodeNotImplemented ErrorCode = "not_implemented"
	ErrorCodeMaintainance   ErrorCode = "undergoing_maintenance"
	ErrorCodeNotFound       ErrorCode = "not_found"
	ErrorCodeInvalidInput   ErrorCode = "invalid_input"
)

func ErrSomethingWentWrong(err error) *APIError {
	return &APIError{
		Message:    "Something went wrong",
		Err:        err,
		StatusCode: http.StatusInternalServerError,
		AppCode:    ErrorCodeInternalError,
		DataError:  false,
	}
}

var (
	ErrNotImplemented = &APIError{
		Message:    "Feature not implemented",
		StatusCode: http.StatusNotFound,
		AppCode:    "not_implemented",
		DataError:  false,
	}
	ErrAPIMaintainance = &APIError{
		Message:    "Maintenance is being done on the API",
		StatusCode: http.StatusServiceUnavailable,
		AppCode:    ErrorCodeMaintainance,
		DataError:  false,
	}
	ErrNotFound = &APIError{
		Message:    "Not found",
		StatusCode: http.StatusNotFound,
		AppCode:    ErrorCodeNotFound,
		DataError:  false,
	}
	ErrInvalidInput = &APIError{
		Message:    "Invalid input",
		StatusCode: http.StatusBadRequest,
		AppCode:    ErrorCodeInvalidInput,
		DataError:  true,
	}

	ErrNilValue               = errors.New("Cannot use nil value")
	ErrUnknownFilterType      = errors.New("Unknown filter type")
	ErrUnsupportedFilterType  = errors.New("Unsupported filter type")
	ErrUnsupportedFilterValue = errors.New("Unsupported filter value")
)
