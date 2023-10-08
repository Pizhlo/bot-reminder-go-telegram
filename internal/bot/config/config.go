package config

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	DBUser    string `mapstructure:"DB_USER,required"`
	DBPass    string `mapstructure:"DB_PASS,required"`
	DBName    string `mapstructure:"DB_NAME,required"`
	DBHost    string `mapstructure:"DB_HOST,required"`
	DBPort    string `mapstructure:"DB_PORT,required"`
	DBAddress string
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	conf := Config{}

	err := viper.ReadInConfig()
	if err != nil {
		return &conf, errors.Wrap(err, `unable to load config`)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		return &conf, errors.Wrap(err, `unable to unmarshal config`)
	}

	conf.DBAddress = fmt.Sprintf(`postgresql://%s:%s@%s:%s/%s?sslmode=disable`, conf.DBUser, conf.DBPass, conf.DBHost, conf.DBPort, conf.DBName)

	fmt.Println(conf)

	return &conf, nil
}
