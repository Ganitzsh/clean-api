package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentAPI struct {
	Config *APIConfig
	Router *gin.Engine
}

func NewPaymentAPI() *PaymentAPI {
	return &PaymentAPI{
		Config: NewAPIConfig(),
	}
}

func (api *PaymentAPI) InitRouter() *PaymentAPI {
	if api.Router != nil {
		return api
	}
	if !api.Config.DevMode {
		gin.SetMode(gin.ReleaseMode)
	}
	api.Router = gin.New()
	return api
}

func (api *PaymentAPI) Start() error {
	api.InitRouter()
	if api.Config.GetTLS() {
		return http.ListenAndServeTLS(
			api.Config.GetHostURL(),
			api.Config.TLSCert,
			api.Config.TLSKey,
			api.Router,
		)
	} else {
		return http.ListenAndServe(
			api.Config.GetFullHost(),
			api.Router,
		)
	}
}
