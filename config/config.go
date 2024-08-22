package config

import (
	"fmt"
	"os"

	commonConfig "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/config"

	"github.com/spf13/viper"
)

var (
	Cold ColdFlat
	Hot  HotFlat
)

type ColdFlat struct {
	AppEnv           string `json:"appEnv" yaml:"appEnv"`
	AppName          string `json:"appName" yaml:"appName"`
	AppPort          uint   `json:"appPort" yaml:"appPort"`
	AppTimezone      string `json:"appTimezone" yaml:"appTimezone"`
	AppApiKey        string `json:"appApiKey" yaml:"appApiKey"`
	AppSecretKey     string `json:"appSecretKey" yaml:"appSecretKey"`
	AppValidationKey string `json:"appValidationKey" yaml:"appValidationKey"`
	JwtSecretKey     string `json:"jwtSecretKey" yaml:"jwtSecretKey"`
	AppHost          string `json:"appHost" yaml:"appHost"`
	// Sentry

	// Postgres Master
	DBMysqlMasterDriver                string `json:"dbMysqlMasterDriver" yaml:"dbMysqlMasterDriver"`
	DBMysqlMasterHost                  string `json:"dbMysqlMasterHost" yaml:"dbMysqlMasterHost"`
	DBMysqlMasterPort                  int    `json:"dbMysqlMasterPort" yaml:"dbMysqlMasterPort"`
	DBMysqlMasterDBName                string `json:"dbMysqlMasterDbName" yaml:"dbMysqlMasterDbName"`
	DBMysqlMasterUser                  string `json:"dbMysqlMasterUser" yaml:"dbMysqlMasterUser"`
	DBMysqlMasterPassword              string `json:"dbMysqlMasterPassword" yaml:"dbMysqlMasterPassword"`
	DBMysqlMasterSSLMode               string `json:"dbMysqlMasterSslMode" yaml:"dbMysqlMasterSslMode"`
	DBMysqlMasterMaxOpenConnections    int    `json:"dbMysqlMasterMaxOpenConnections" yaml:"dbMysqlMasterMaxOpenConnections"`
	DBMysqlMasterMaxLifeTimeConnection int    `json:"dbMysqlMasterMaxLifeTimeConnection" yaml:"dbMysqlMasterMaxLifeTimeConnection"`
	DBMysqlMasterMaxIdleConnections    int    `json:"dbMysqlMasterMaxIdleConnections" yaml:"dbMysqlMasterMaxIdleConnections"`
	DBMysqlMasterMaxIdleTimeConnection int    `json:"dbMysqlMasterMaxIdleTimeConnection" yaml:"dbMysqlMasterMaxIdleTimeConnection"`

	// Postgres Slave
	DBMysqlSlaveDriver                string `json:"dbMysqlSlaveDriver" yaml:"dbMysqlSlaveDriver"`
	DBMysqlSlaveHost                  string `json:"dbMysqlSlaveHost" yaml:"dbMysqlSlaveHost"`
	DBMysqlSlavePort                  int    `json:"dbMysqlSlavePort" yaml:"dbMysqlSlavePort"`
	DBMysqlSlaveDBName                string `json:"dbMysqlSlaveDbName" yaml:"dbMysqlSlaveDbName"`
	DBMysqlSlaveUser                  string `json:"dbMysqlSlaveUser" yaml:"dbMysqlSlaveUser"`
	DBMysqlSlavePassword              string `json:"dbMysqlSlavePassword" yaml:"dbMysqlSlavePassword"`
	DBMysqlSlaveSSLMode               string `json:"dbMysqlSlaveSslMode" yaml:"dbMysqlSlaveSslMode"`
	DBMysqlSlaveMaxOpenConnections    int    `json:"dbMysqlSlaveMaxOpenConnections" yaml:"dbMysqlSlaveMaxOpenConnections"`
	DBMysqlSlaveMaxLifeTimeConnection int    `json:"dbMysqlSlaveMaxLifeTimeConnection" yaml:"dbMysqlSlaveMaxLifeTimeConnection"`
	DBMysqlSlaveMaxIdleConnections    int    `json:"dbMysqlSlaveMaxIdleConnections" yaml:"dbMysqlSlaveMaxIdleConnections"`
	DBMysqlSlaveMaxIdleTimeConnection int    `json:"dbMysqlSlaveMaxIdleTimeConnection" yaml:"dbMysqlSlaveMaxIdleTimeConnection"`

	// Redis
	RedisHost     string `json:"redisHost" yaml:"redisHost"`
	RedisUsername string `json:"redisUsername" yaml:"redisUsername"`
	RedisPassword string `json:"redisPassword" yaml:"redisPassword"`
}

type HotFlat struct {
	AppDebug              bool   `json:"appDebug" yaml:"appDebug"`
	LoggerDebug           bool   `json:"loggerDebug" yaml:"loggerDebug"`
	SentryDebug           bool   `json:"sentryDebug" yaml:"sentryDebug"`
	ShutDownDelayInSecond uint64 `json:"shutDownDelayInSecond" yaml:"shutDownDelayInSecond"`
}

func Init() {
	consulURL := os.Getenv("CONSUL_HTTP_URL")

	err := commonConfig.BindFromFile(&Cold, "config.cold.json", ".")
	if err != nil {
		fmt.Printf("failed load cold config from file: %s", viper.ConfigFileUsed())
		err = commonConfig.BindFromConsul(
			&Cold,
			consulURL,
			fmt.Sprintf("%s/%s", os.Getenv("CONSUL_HTTP_KEY"), "cold"),
		)
		if err != nil {
			panic(err)
		}
	}

	err = commonConfig.BindFromFile(&Hot, "config.hot.json", ".")
	if err != nil {
		fmt.Printf("failed load hot config from file: %s", viper.ConfigFileUsed())

		interval, err := commonConfig.LoadConsulIntervalFromEnv()
		if err != nil {
			panic(err)
		}

		err = commonConfig.BindAndWatchFromConsul(
			&Hot,
			consulURL,
			fmt.Sprintf("%s/%s", os.Getenv("CONSUL_HTTP_KEY"), "hot"),
			interval,
		)
		if err != nil {
			panic(err)
		}
	}
}
