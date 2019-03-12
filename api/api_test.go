package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ganitzsh/f3-te/api"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func doHTTPReq(handler http.Handler, method string, url string, body url.Values) *http.Response {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	handler.ServeHTTP(rr, req)
	return rr.Result()
}

func TestNotFound(t *testing.T) {
	handler := api.Routes()
	resp := doHTTPReq(handler, http.MethodGet, "/unknown", nil)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	d := api.JSENDData{}
	if !assert.NoError(t, json.Unmarshal(body, &d)) {
		t.FailNow()
	}
}

func testListPayments(db *mockDB) func(*testing.T) {
	api.InitStore(db.Store)
	handler := api.Routes()
	return func(t *testing.T) {
		resp := doHTTPReq(handler, http.MethodGet, "/payments", nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		d := api.JSENDData{Data: []*api.Payment{}}
		if !assert.NoError(t, json.Unmarshal(body, &d)) {
			t.FailNow()
		}
		assert.Len(t, d.Data, 3)
	}
}

func testGetPayment(db *mockDB) func(*testing.T) {
	api.InitStore(db.Store)
	handler := api.Routes()
	return func(t *testing.T) {
		resp := doHTTPReq(handler, http.MethodGet, "/payments/unknown", nil)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, api.ContentTypeJSON, resp.Header.Get(api.HeaderContentType))
		resp = doHTTPReq(handler, http.MethodGet, "/payments/"+uuid.New().String(), nil)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, api.ContentTypeJSON, resp.Header.Get(api.HeaderContentType))
		resp = doHTTPReq(handler, http.MethodGet, "/payments/"+db.ID1.String(), nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, api.ContentTypeJSON, resp.Header.Get(api.HeaderContentType))
		body, _ := ioutil.ReadAll(resp.Body)
		d := api.JSENDData{Data: new(api.Payment)}
		if !assert.NoError(t, json.Unmarshal(body, &d)) {
			t.FailNow()
		}
		assert.NotNil(t, d.Data)
	}
}

func testSavePayment(db *mockDB) func(*testing.T) {
	api.InitStore(db.Store)
	return func(t *testing.T) {
	}
}

func TestListPaymentsWithInMemStore(t *testing.T) {
	testListPayments(newTestDBInMem())(t)
}

func TestSavePaymentWithInMemStore(t *testing.T) {
	testSavePayment(newTestDBInMem())(t)
}

func TestGetPaymentWithInMemStore(t *testing.T) {
	testGetPayment(newTestDBInMem())(t)
}
