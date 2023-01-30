package config

import "github.com/spf13/viper"

type Config struct {
	MysqlUser         string `mapstructure:"MYSQL_USER"`
	MysqlPassword     string `mapstructure:"MYSQL_PASSWORD"`
	MysqlRootPassword string `mapstructure:"MYSQL_ROOT_PASSWORD"`
	MysqlDatabase     string `mapstructure:"MYSQL_DATABASE"`
	Addr              string `mapstructure:"ADDR"`
	HashSalt          string `mapstructure:"HASH_SALT"`
	SigningKey        string `mapstructure:"SIGNING_KEY"`
	TokenTtl          int    `mapstructure:"TOKEN_TTL"`
}

func InitConfig() (config Config, err error) {
	viper.AddConfigPath("./interal/config")
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
