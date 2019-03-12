package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

type APIError struct {
	DataError  bool   `json:"-"`
	Message    string `json:"error"`
	StatusCode int    `json:"-"`
	AppCode    string `json:"code,omitempty"`
	Err        error  `json:"-"`
}

func (e APIError) Error() string {
	return e.Message
}

func (e *APIError) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	status := "error"
	if e.DataError {
		status = "fail"
	}
	payload := JSENDData{
		Code:   e.StatusCode,
		Data:   e,
		Status: status,
	}
	render.JSON(w, r, payload)
	return nil
}

func NewAPIError(dataError bool, message string) *APIError {
	return &APIError{
		DataError: dataError,
		Message:   message,
	}
}

func ErrSomethingWentWrong(err error) *APIError {
	return &APIError{
		Message:    "Something went wrong",
		Err:        err,
		StatusCode: http.StatusInternalServerError,
		AppCode:    "internal_error",
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
	ErrNotFound = &APIError{
		Message:    "Not found",
		StatusCode: http.StatusNotFound,
		AppCode:    "not_found",
		DataError:  false,
	}
	ErrInvalidInput = &APIError{
		Message:    "Invalid input",
		StatusCode: http.StatusBadRequest,
		AppCode:    "invalid_input",
		DataError:  true,
	}

	ErrNilValue               = errors.New("Cannot use nil value")
	ErrUnknownFilterType      = errors.New("Unknown filter type")
	ErrUnsupportedFilterType  = errors.New("Unsupported filter type")
	ErrUnsupportedFilterValue = errors.New("Unsupported filter value")
)
