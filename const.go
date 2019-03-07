package main

const (
	ReleaseName = "elliot"
	Version     = "0.0.1"
	EmptyString = ""

	DefaultNodeName = "Payment API"
	DefaultHost     = "127.0.0.1"
	DefaultPort     = "8080"
	DefaultDevMode  = true
	DefaultTLS      = false
	DefaultTLSKey   = "private_key"
	DefaultTLSCert  = "cert"

	EnvPrefix            = "api"
	ConfigFileName       = "config"
	ConfigKeyHost        = "host"
	ConfigKeyPort        = "port"
	ConfigKeyTLS         = "tls.enabled"
	ConfigKeyTLSKey      = "tls.key"
	ConfigKeyTLSCert     = "tls.cert"
	ConfigKeyCORSOrigins = "cors.origins"
	ConfigKeyCORSMethods = "cors.methods"
	ConfigKeyCORSHeaders = "cors.headers"
	ConfigKeyDevMode     = "dev_mode"
	ConfigKeyNodeName    = "name"

	PaymentIDPrefix = "payment_id"
)
