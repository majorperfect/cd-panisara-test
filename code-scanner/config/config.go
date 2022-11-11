package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Host   string `mapstructure:"server_host"`
	Port   string `mapstructure:"server_port"`
	DBHost string `mapstructure:"db_host"`
	DBPort string `mapstructure:"db_port"`
}

func SetupConfig() *Config {
	viper.SetConfigFile("./config/.env")
	_ = viper.ReadInConfig()

	// environment variable over .env file
	viper.AutomaticEnv()

	cfg := &Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panic("unable to decode into config struct, %v \r\n", err)
	}

	///////// BINDING TO GO ENV FOR FUTURE USES /////////
	os.Setenv("ENV", viper.GetString("ENV"))

	os.Setenv("DB_HOST", cfg.DBHost)
	os.Setenv("DB_PORT", cfg.DBPort)
	os.Setenv("SERVER_HOST", cfg.Host)
	os.Setenv("SERVER_PORT", cfg.Port)

	// TLS certificate
	os.Setenv("TLS_CERT_FILE", viper.GetString("TLS_CERT_FILE"))
	os.Setenv("TLS_CLIENT_CERT_FILE", viper.GetString("TLS_CLIENT_CERT_FILE"))
	os.Setenv("TLS_CLIENT_CERT_KEY_FILE", viper.GetString("TLS_CLIENT_CERT_KEY_FILE"))
	////////////////////////////////////////////////////
	// os.Setenv("JWT_KEY", viper.GetString("JWT_KEY"))
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		fmt.Println(variable[0], "=>", variable[1])
	}

	return cfg
}
