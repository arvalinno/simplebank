package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// LoadConfig reads configuration from file or env variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // can be json, xml

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		// If config file is not found, read from environment variables
		config.DBDriver = viper.GetString("DB_DRIVER")
		config.DBSource = viper.GetString("DB_SOURCE")
		config.ServerAddress = viper.GetString("SERVER_ADDRESS")
		config.TokenSymmetricKey = viper.GetString("TOKEN_SYMMETRIC_KEY")
		config.AccessTokenDuration = viper.GetDuration("ACCESS_TOKEN_DURATION")
		config.RefreshTokenDuration = viper.GetDuration("REFRESH_TOKEN_DURATION")

		// Check if essential configurations are still missing
		if config.DBSource == "" {
			return config, err
		}

		fmt.Println("debug env:", config, err)
	} else {
		err = viper.Unmarshal(&config)
	}

	fmt.Println("debug env2:", config, err)

	return config, err
}
