package api

import (
	"net/http"

	"github.com/go-chi/render"
)

const (
	JSENDDataStatusSuccess = "success"
	JSENDDataStatusFail    = "fail"
	JSENDDataStatusError   = "error"
)

type JSENDData struct {
	Data   interface{} `json:"data"`
	Code   int         `json:"code"`
	Status string      `json:"status"`
}

func NewJSENDData(data interface{}, code ...int) *JSENDData {
	var overrideCode int
	if code != nil {
		overrideCode = code[0]
	}
	return &JSENDData{
		Data:   data,
		Code:   overrideCode,
		Status: JSENDDataStatusSuccess,
	}
}

func (p *JSENDData) Render(w http.ResponseWriter, r *http.Request) error {
	status := JSENDDataStatusSuccess
	code := http.StatusOK
	if apiError, ok := p.Data.(*APIError); ok {
		status = JSENDDataStatusError
		code = http.StatusInternalServerError
		if apiError.DataError {
			status = JSENDDataStatusFail
		}
		if apiError.StatusCode != 0 {
			code = apiError.StatusCode
		}
		if p.Code != 0 {
			code = p.Code
		}
		p.Code = code
		p.Status = status
		render.Status(r, p.Code)
		return nil
	}
	p.Code = code
	p.Status = status
	render.Status(r, p.Code)
	return nil
}
