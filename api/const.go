package api

const (
	ReleaseName       = "elliot"
	Version           = "0.0.1"
	CurrentAPIVersion = APIVersion(1)

	URLRoot = "/"

	APIV1ContentTypes = "application/json,application/json+v1"
	APIV1Prefix       = "/v1"

	ReqDataKey = "data"
	ReqCodeKey = "code"

	DatabaseTypeInMem = "inmem"
	DatabaseTypeMongo = "mongo"

	DefaultNodeName        = "Payment API"
	DefaultHost            = "127.0.0.1"
	DefaultPort            = "8080"
	DefaultDevMode         = true
	DefaultTLS             = false
	DefaultTLSKey          = "private_key"
	DefaultTLSCert         = "cert"
	DefaultMongoDatabase   = "payment_api"
	DefaultMongoCollection = "payments"
	DefaultMongoURI        = "localhost"
	DefaultMongoMaxRetries = 10
	DefaultDBType          = DatabaseTypeInMem

	EnvPrefix                = "api"
	ConfigFileName           = "config"
	ConfigKeyHost            = "host"
	ConfigKeyPort            = "port"
	ConfigKeyTLS             = "tls.enabled"
	ConfigKeyTLSKey          = "tls.key"
	ConfigKeyTLSCert         = "tls.cert"
	ConfigKeyCORSOrigins     = "cors.origins"
	ConfigKeyCORSMethods     = "cors.methods"
	ConfigKeyCORSHeaders     = "cors.headers"
	ConfigKeyDatabaseType    = "database.type"
	ConfigKeyMongoDatabase   = "database.mongo.database"
	ConfigKeyMongoCollection = "database.mongo.collection"
	ConfigKeyMongoURI        = "database.mongo.uri"
	ConfigKeyMongoMaxRetries = "database.mongo.max_retries"
	ConfigKeyDevMode         = "dev_mode"
	ConfigKeyNodeName        = "name"

	PaymentIDPrefix = "payment_id"

	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json; charset=utf-8"
)
