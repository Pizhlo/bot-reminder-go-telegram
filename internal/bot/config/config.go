package config

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	DBUser             string        `mapstructure:"POSTGRES_USER,required"`
	DBPass             string        `mapstructure:"POSTGRES_PASSWORD,required"`
	DBName             string        `mapstructure:"POSTGRES_DB,required"`
	DBHost             string        `mapstructure:"POSTGRES_HOST,required"`
	DBPort             string        `mapstructure:"POSTGRES_PORT,required"`
	Timeout            time.Duration `mapstructure:"TIMEOUT,required"`
	TelegramMaxMsgSize int           `mapstructure:"MAX_MSG_SIZE,required"`
	DBAddress          string
	Token              string `mapstructure:"TOKEN,required"`
	ChannelID          int64  `mapstructure:"CHANNEL_ID,required"`
	BotURL             string `mapstructure:"BOT_URL"`
	LogLvl             string `mapstructure:"LOG_LEVEL"`
}

func LoadConfig(file, path string) (*Config, error) {
	viper.SetConfigFile(file)
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

	return &conf, nil
}
