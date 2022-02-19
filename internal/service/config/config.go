package config

import (
	"strings"

	"github.com/spf13/viper"
)

func Init() error {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("log.level", "trace")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("port", "8080")
	viper.SetDefault("server.host", viper.GetString("host"))
	viper.SetDefault("server.port", viper.GetString("port"))
	viper.SetDefault("database.url", "postgresql://postgres:postgres@localhost/go_base")

	return nil
}
