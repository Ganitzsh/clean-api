package api

import (
	"net/http"
	"strings"

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
	viper.AutomaticEnv()
}

func InitStore(s PaymentStore) {
	store = s
}
