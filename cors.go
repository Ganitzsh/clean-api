package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type CORSSettings struct {
	AllowedOrigins []string `json:"allowed_origins"`
	AllowedMethods []string `json:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers"`
}

func NewCORSSettings() *CORSSettings {
	return &CORSSettings{
		AllowedHeaders: viper.GetStringSlice(ConfigKeyCORSHeaders),
		AllowedMethods: viper.GetStringSlice(ConfigKeyCORSMethods),
		AllowedOrigins: viper.GetStringSlice(ConfigKeyCORSOrigins),
	}
}

const (
	AccessControlAllowOrigin  = "Access-Control-Allow-Origin"
	AccessControlAllowMethods = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders = "Access-Control-Allow-Headers"
)

func CORS(settings CORSSettings) func(*gin.Context) {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Header(AccessControlAllowOrigin, strings.Join(settings.AllowedOrigins, ","))
			c.Header(AccessControlAllowMethods, strings.Join(settings.AllowedMethods, ","))
			c.Header(AccessControlAllowHeaders, strings.Join(settings.AllowedHeaders, ","))
		}
	}
}
