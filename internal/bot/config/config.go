package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

	logrus.Debug("hello world")

	conf := Config{}

	err := viper.ReadInConfig()
	if err != nil {
		return loadEnvVariables()
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		return &conf, errors.Wrap(err, `unable to unmarshal config`)
	}

	conf.DBAddress = fmt.Sprintf(`postgresql://%s:%s@%s:%s/%s?sslmode=disable`, conf.DBUser, conf.DBPass, conf.DBHost, conf.DBPort, conf.DBName)

	return &conf, nil
}

func loadEnvVariables() (*Config, error) {
	conf := &Config{}

	dbUser := os.Getenv("POSTGRES_USER")
	if len(dbUser) == 0 {
		return nil, errors.New("POSTGRES_USER is not set")
	}

	conf.DBUser = dbUser

	dbPass := os.Getenv("POSTGRES_PASSWORD")
	if len(dbPass) == 0 {
		return nil, errors.New("POSTGRES_PASSWORD is not set")
	}

	conf.DBPass = dbPass

	dbName := os.Getenv("POSTGRES_DB")
	if len(dbName) == 0 {
		return nil, errors.New("POSTGRES_DB is not set")
	}

	conf.DBName = dbName

	dbHost := os.Getenv("POSTGRES_HOST")
	if len(dbHost) == 0 {
		return nil, errors.New("POSTGRES_HOST is not set")
	}

	conf.DBHost = dbHost

	dbPort := os.Getenv("POSTGRES_PORT")
	if len(dbPort) == 0 {
		return nil, errors.New("POSTGRES_PORT is not set")
	}

	conf.DBPort = dbPort

	token := os.Getenv("TOKEN")
	if len(token) == 0 {
		return nil, errors.New("TOKEN is not set")
	}

	conf.Token = token

	timeout := os.Getenv("TIMEOUT")
	if len(timeout) == 0 {
		return nil, errors.New("TIMEOUT is not set")
	}

	d, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, fmt.Errorf("unable to parse duration %s: %v", timeout, err)
	}

	conf.Timeout = d

	maxMsgSizeStr := os.Getenv("MAX_MSG_SIZE")
	if len(maxMsgSizeStr) == 0 {
		return nil, errors.New("MAX_MSG_SIZE is not set")
	}

	size, err := strconv.Atoi(maxMsgSizeStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse maxMsgSize %s: %v", maxMsgSizeStr, err)
	}

	conf.TelegramMaxMsgSize = size

	channelStr := os.Getenv("CHANNEL_ID")
	if len(channelStr) == 0 {
		return nil, errors.New("CHANNEL_ID is not set")
	}

	channel, err := strconv.Atoi(channelStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse channelID %s: %v", channelStr, err)
	}

	conf.ChannelID = int64(channel)

	logLvl := os.Getenv("LOG_LEVEL")
	if len(logLvl) == 0 {
		return nil, errors.New("LOG_LEVEL is not set")
	}

	conf.LogLvl = logLvl

	return conf, nil
}
