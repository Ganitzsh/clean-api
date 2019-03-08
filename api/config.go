package api

import (
	"fmt"

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

type APIConfig struct {
	NodeName string        `json:"node_name"`
	DevMode  bool          `json:"dev_mode"`
	Host     string        `json:"host" validate:"required"`
	Port     string        `json:"port" validate:"required"`
	TLS      bool          `json:"tls"`
	TLSKey   string        `json:"tls_key"`
	TLSCert  string        `json:"tls_cert"`
	Cors     *CORSSettings `json:"cors"`
}

// NewAPIConfig creates a new APIConfig struct.
func NewAPIConfig() *APIConfig {
	return &APIConfig{
		NodeName: viper.GetString(ConfigKeyNodeName),
		DevMode:  viper.GetBool(ConfigKeyDevMode),
		Host:     viper.GetString(ConfigKeyHost),
		Port:     viper.GetString(ConfigKeyPort),
		TLS:      viper.GetBool(ConfigKeyTLS),
		TLSKey:   viper.GetString(ConfigKeyTLSKey),
		TLSCert:  viper.GetString(ConfigKeyTLSCert),
		Cors:     NewCORSSettings(),
	}
}

// SetHost sets Host to value
func (c *APIConfig) SetHost(value string) *APIConfig {
	c.Host = value
	return c
}

// SetHost sets Port to value
func (c *APIConfig) SetPort(value string) *APIConfig {
	c.Port = value
	return c
}

// SetTLS sets TLS to value
func (c *APIConfig) SetTLS(value bool) *APIConfig {
	c.TLS = value
	return c
}

// SetDevMode sets DevMode to value
func (c *APIConfig) SetDevMode(value bool) *APIConfig {
	c.DevMode = value
	return c
}

// GetHost returns the value of Host
func (c *APIConfig) GetHost() string {
	return c.Host
}

// GetTLS returns the value of TLS
func (c *APIConfig) GetTLS() bool {
	return c.TLS
}

// GetPort returns the value of Port
func (c *APIConfig) GetPort() string {
	return c.Port
}

// GetHostURL generates a string representing the full host of the API.
// When the Host or Port is an empty string it will use the default.
func (c *APIConfig) GetHostURL() string {
	proto := "http"
	if c.TLS {
		proto = "https"
	}
	host := c.Host
	if host == "" {
		host = DefaultHost
	}
	port := c.Port
	if port == "" {
		port = DefaultPort
	}
	return fmt.Sprintf("%s://%s:%s", proto, host, port)
}

// GetHostURL generates a string representing the full host of the API.
// When the Host or Port is an empty string it will use the default.
func (c *APIConfig) GetFullHost() string {
	host := c.Host
	if host == "" {
		host = DefaultHost
	}
	port := c.Port
	if port == "" {
		port = DefaultPort
	}
	return fmt.Sprintf("%s:%s", host, port)
}
