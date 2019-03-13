package api_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ganitzsh/f3-te/api"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	api.InitConfig()
}

func readErrorCode(body []byte) api.ErrorCode {
	if body == nil || len(body) == 0 {
		return ""
	}
	apiErr := &api.APIError{}
	tmp := api.JSENDData{Data: apiErr}
	if err := json.Unmarshal(body, &tmp); err == nil {
		return apiErr.AppCode
	}
	return ""
}

func doHTTPReq(handler http.Handler, method string, url string, body string) *http.Response {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, bytes.NewBufferString(body))
	req.Header.Add(api.HeaderContentType, "application/json")
	handler.ServeHTTP(rr, req)
	return rr.Result()
}

func TestNotFound(t *testing.T) {
	handler := api.Routes()
	resp := doHTTPReq(handler, http.MethodGet, "/unknown", "")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, api.ErrNotFound.AppCode, readErrorCode(body))
	d := api.JSENDData{}
	if !assert.NoError(t, json.Unmarshal(body, &d)) {
		t.FailNow()
	}
}

func TestWrongContentType(t *testing.T) {
	handler := api.Routes()
	req, _ := http.NewRequest(http.MethodPost, "/v1/payments", bytes.NewBufferString(""))
	rr := httptest.NewRecorder()
	req.Header.Add(api.HeaderContentType, "some/madeup-ct")
	handler.ServeHTTP(rr, req)
	resp := rr.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, api.ErrInvalidInput.AppCode, readErrorCode(body))
}

func testListPayments(db *mockDB) func(*testing.T) {
	api.SetStore(db.Store)
	handler := api.Routes()
	return func(t *testing.T) {
		resp := doHTTPReq(handler, http.MethodGet, "/v1/payments", "")
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
	api.SetStore(db.Store)
	handler := api.Routes()
	return func(t *testing.T) {
		resp := doHTTPReq(handler, http.MethodGet, "/v1/payments/unknown", "")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, api.ContentTypeJSON, resp.Header.Get(api.HeaderContentType))
		resp = doHTTPReq(handler, http.MethodGet, "/v1/payments/"+uuid.New().String(), "")
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Equal(t, api.ContentTypeJSON, resp.Header.Get(api.HeaderContentType))
		resp = doHTTPReq(handler, http.MethodGet, "/v1/payments/"+db.ID1.String(), "")
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
	api.SetStore(db.Store)
	return func(t *testing.T) {
		handler := api.Routes()
		resp := doHTTPReq(handler, http.MethodPost, "/v1/payments", "")
		body, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, api.ErrInvalidInput.AppCode, readErrorCode(body))

		payment := newMockPayment()
		payment.Amount = "42"
		b, _ := json.Marshal(payment)
		resp = doHTTPReq(handler, http.MethodPost, "/v1/payments", string(b))
		body, _ = ioutil.ReadAll(resp.Body)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, api.ErrorCode(""), readErrorCode(body))
		p := api.JSENDData{Data: new(api.Payment)}
		if !assert.NoError(t, json.Unmarshal(body, &p)) {
			t.FailNow()
		}
		respPayment := p.Data.(*api.Payment)
		assert.NotEqual(t, respPayment.ID, "")
		assert.Equal(t, "42", respPayment.Amount)
		id := respPayment.ID
		ref := respPayment.EndToEndReference

		payment.Amount = "84"
		b, _ = json.Marshal(payment)
		resp = doHTTPReq(handler, http.MethodPost, "/v1/payments/"+id.String(), string(b))
		body, _ = ioutil.ReadAll(resp.Body)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, api.ErrorCode(""), readErrorCode(body))
		p = api.JSENDData{Data: new(api.Payment)}
		if !assert.NoError(t, json.Unmarshal(body, &p)) {
			t.FailNow()
		}
		assert.Equal(t, id, p.Data.(*api.Payment).ID)
		assert.Equal(t, ref, p.Data.(*api.Payment).EndToEndReference)
		assert.Equal(t, "84", p.Data.(*api.Payment).Amount)
	}
}

func testDeletePayment(db *mockDB) func(*testing.T) {
	api.SetStore(db.Store)
	handler := api.Routes()
	return func(t *testing.T) {
		resp := doHTTPReq(handler, http.MethodDelete, "/v1/payments/"+uuid.New().String(), "")
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		resp = doHTTPReq(handler, http.MethodDelete, "/v1/payments/"+db.ID1.String(), "")
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		resp = doHTTPReq(handler, http.MethodGet, "/v1/payments/"+db.ID1.String(), "")
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	}
}

func TestListPaymentsWithInMemStore(t *testing.T) {
	testListPayments(newTestDBInMem())(t)
}

func TestSavePaymentWithInMemStore(t *testing.T) {
	testSavePayment(newTestDBInMem())(t)
}

func TestDeletePaymentWithInMemStore(t *testing.T) {
	testDeletePayment(newTestDBInMem())(t)
}

func TestGetPaymentWithInMemStore(t *testing.T) {
	testGetPayment(newTestDBInMem())(t)
}
