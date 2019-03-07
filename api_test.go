package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPayments(t *testing.T) {
	api := NewPaymentAPI().InitRouter()
	recorder := httptest.NewRecorder()
	url := api.Config.GetHostURL() + "/payments"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	api.Router.ServeHTTP(recorder, req)
}
