package api

const (
	ReleaseName       = "elliot"
	Version           = "0.0.1"
	CurrentAPIVersion = APIVersion(1)

	APIV1ContentTypes = "application/json,application/json+v1"
	APIV1Prefix       = "/v1"

	ReqDataKey = "data"
	ReqCodeKey = "code"

	DefaultNodeName = "Document API"
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

	DocumentIDPrefix = "payment_id"
)
