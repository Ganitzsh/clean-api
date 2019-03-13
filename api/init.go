package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName(ConfigFileName)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Warnf("could not read config: %v", err)
	}
	viper.SetDefault(ConfigKeyHost, DefaultHost)
	viper.SetDefault(ConfigKeyPort, DefaultPort)
	viper.SetDefault(ConfigKeyDevMode, DefaultDevMode)
	viper.SetDefault(ConfigKeyTLS, DefaultTLS)
	viper.SetDefault(ConfigKeyNodeName, DefaultNodeName)
	viper.SetDefault(ConfigKeyCORSHeaders, []string{"*"})
	viper.SetDefault(ConfigKeyCORSMethods, []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
	})
	viper.SetDefault(ConfigKeyCORSOrigins, []string{"*"})
	viper.SetDefault(ConfigKeyMongoDatabase, DefaultMongoDatabase)
	viper.SetDefault(ConfigKeyMongoCollection, DefaultMongoCollection)
	viper.SetDefault(ConfigKeyMongoURI, DefaultMongoURI)
	viper.SetDefault(ConfigKeyMongoMaxRetries, DefaultMongoMaxRetries)
	viper.SetDefault(ConfigKeyDatabaseType, DatabaseTypeInMem)
	viper.AutomaticEnv()
	config = NewAPIConfig()
}

func mongoHealthCheck() {
	var attempts int
	for {
		prevState := mongoHealthy
		if mongo == nil {
			logrus.Error("No mongo handler")
			mongoHealthy = false
		} else {
			if err := mongo.Ping(); err != nil {
				mongoHealthy = false
				attempts += 1
				logrus.Errorf(
					"Mongo: could not ping database: %v (%d of %d)",
					err, attempts, config.Mongo.MaxRetries,
				)
				if attempts >= config.Mongo.MaxRetries {
					logrus.Fatalf("Could not reach database afer %d attempts", attempts)
				}
			} else {
				mongoHealthy = true
			}
			if prevState == false && mongoHealthy {
				logrus.Info("Mongo: connection reestablished")
				attempts = 0
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func getMongoCollection() (*MgoWrapCollection, error) {
	var err error
	logrus.Info("Mongo: Connecting")
	mongo, err = mgo.DialWithTimeout(config.Mongo.URI, 1*time.Second)
	if err != nil {
		return nil, err
	}
	if err = mongo.Ping(); err != nil {
		return nil, errors.New("Could not ping the database")
	}
	db := mongo.DB(config.Mongo.Database)
	if _, err = db.CollectionNames(); err != nil {
		return nil, errors.New("Could not retrieve collections, are you logged in?")
	}
	mongoHealthy = true
	go mongoHealthCheck()
	return &MgoWrapCollection{db.C(config.Mongo.Collection)}, nil
}

func initMongo() {
	var attempts int
	for {
		c, err := getMongoCollection()
		if err != nil {
			attempts += 1
			logrus.Errorf(
				"Mongo: could not connect %v (attempt %d on %d)",
				err, attempts, config.Mongo.MaxRetries,
			)
			if attempts >= config.Mongo.MaxRetries {
				logrus.Fatalf("Could not reach database afer %d attempts", attempts)
			}
		} else {
			store = NewPaymentMongoStore(c)
			break
		}
	}
}

func InitStore() {
	if config == nil {
		logrus.Fatal("No configuration found")
	}
	switch config.DBType {
	case DatabaseTypeInMem:
		store = NewPaymentInMemStore()
		break
	case DatabaseTypeMongo:
		initMongo()
		break
	default:
		logrus.Fatal("Unknown or empty database type")
	}
}

func SetStore(s PaymentStore) {
	store = s
}
