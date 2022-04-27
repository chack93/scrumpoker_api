package config

import (
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

func Init() error {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("log.level", "trace")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", "8080")
	viper.SetDefault("server.host", viper.GetString("host"))
	viper.SetDefault("server.port", viper.GetString("port"))
	viper.SetDefault("database.url", "postgresql://postgres:postgres@localhost/scrumpoker_api")
	viper.SetDefault("msgqueue.nats.url", nats.DefaultURL)

	return nil
}
