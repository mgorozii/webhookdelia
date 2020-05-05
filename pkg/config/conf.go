package config

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/juju/errors"
	"github.com/spf13/viper"
)

type Conf struct {
	TelegramToken string
	PublicURL     string
	TelegramURL   string
	Port          string
	Release       bool
	StoreOptions  map[string]interface{}
}

func (c Conf) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.TelegramToken, validation.Required, validation.Length(1, 0)),
		validation.Field(&c.PublicURL, validation.Required, validation.Length(1, 0)),
	)
}

func New() (Conf, error) {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/webhookdelia")

	viper.SetDefault("telegram_url", "https://api.telegram.org")
	viper.SetDefault("release", false)
	viper.SetDefault("port", "8080")

	viper.SetDefault("store_file_directory", "/tmp/webhookdelia")

	viper.SetDefault("store_postgres_connection_url", "postgres://postgres@/webhookdelia?sslmode=disable")
	viper.SetDefault("store_postgres_table_name", "webhookdelia")
	viper.SetDefault("store_postgres_max_open_connections", 100)
	viper.SetDefault("store_redis_address", "localhost:6379")
	viper.SetDefault("store_redis_password", "")
	viper.SetDefault("store_redis_db", 0)

	// try to read the conf file
	_ = viper.ReadInConfig()

	// and then use env
	viper.AutomaticEnv()

	conf := Conf{
		TelegramToken: viper.GetString("telegram_token"),
		PublicURL:     viper.GetString("public_url"),
		TelegramURL:   viper.GetString("telegram_url"),
		Port:          viper.GetString("port"),
		Release:       viper.GetBool("release"),
		StoreOptions:  map[string]interface{}{},
	}

	for _, storeOption := range viper.AllKeys() {
		storeOption = strings.ToLower(storeOption)
		if !strings.HasPrefix(storeOption, "store_") {
			continue
		}

		conf.StoreOptions[storeOption] = viper.Get(storeOption)
	}

	if err := conf.Validate(); err != nil {
		return conf, errors.Trace(err)
	}

	fmt.Println("conf", conf)

	return conf, nil
}
