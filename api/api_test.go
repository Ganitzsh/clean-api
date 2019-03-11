package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ganitzsh/f3-te/api"
	"github.com/stretchr/testify/assert"
)

func testListPayment(store api.DocumentStore) func(*testing.T) {
	api.InitStore(store)
	handler := api.Routes()
	return func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/payments", nil)
		handler.ServeHTTP(rr, req)
		resp := rr.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		d := api.JSENDData{Data: []*api.Document{}}
		if !assert.NoError(t, json.Unmarshal(body, &d)) {
			t.FailNow()
		}
		assert.Len(t, d.Data, 4)
	}
}

func testGetPayment(store api.DocumentStore) func(*testing.T) {
	api.InitStore(store)
	return func(t *testing.T) {
		handler := http.HandlerFunc(api.GetPayment)
		assert.HTTPError(t, handler, http.MethodGet, "/payments", nil)
	}
}

func testSavePayment(store api.DocumentStore) func(*testing.T) {
	api.InitStore(store)
	return func(t *testing.T) {
	}
}

func TestListPaymentsWithInMemStore(t *testing.T) {
	testListPayment(newTestDBInMem().Store)(t)
}

func TestSavePaymentWithInMemStore(t *testing.T) {
	testSavePayment(newTestDBInMem().Store)(t)
}

func TestGetPaymentWithInMemStore(t *testing.T) {
	testGetPayment(newTestDBInMem().Store)(t)
}
