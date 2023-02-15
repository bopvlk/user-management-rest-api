package config

import (
	"git.foxminded.com.ua/3_REST_API/interal/apperrors"
	"github.com/spf13/viper"
)

type Config struct {
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	Port       string `mapstructure:"PORT"`
	HashSalt   string `mapstructure:"HASH_SALT"`
	SigningKey string `mapstructure:"SIGNING_KEY"`
	TokenTtl   int    `mapstructure:"TOKEN_TTL"`
}

func InitConfig() (config *Config, err error) {
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./build/package/config")
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return nil, apperrors.ConfigReadErr.AppendMessage(err)
	}

	if err = viper.Unmarshal(&config); err != nil {
		return nil, apperrors.ConfigUnmarshallErr.AppendMessage(err)
	}
	return
}
