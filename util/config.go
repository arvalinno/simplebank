package util

import (
	"log"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/spf13/viper"
)

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
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
		if err == err.(viper.ConfigFileNotFoundError) {
			// If config file is not found, read from environment variables
			log.Println("Config file not found, relying on environment variables")

			config.DBDriver = viper.GetString("DB_DRIVER")
			config.DBSource = viper.GetString("DB_SOURCE")
			config.HTTPServerAddress = viper.GetString("HTTP_SERVER_ADDRESS")
			config.GRPCServerAddress = viper.GetString("GRPC_SERVER_ADDRESS")
			config.TokenSymmetricKey = viper.GetString("TOKEN_SYMMETRIC_KEY")
			config.AccessTokenDuration = viper.GetDuration("ACCESS_TOKEN_DURATION")
			config.RefreshTokenDuration = viper.GetDuration("REFRESH_TOKEN_DURATION")

			// load midtrans config
			midtrans.ServerKey = viper.GetString("SERVER_KEY")
			midtrans.Environment = midtrans.Sandbox

			// Check if essential configurations are still missing
			if config.DBSource == "" {
				return config, err
			}

			return config, nil
		} else {
			return config, err
		}

	} else {
		err = viper.Unmarshal(&config)

		// load midtrans config
		midtrans.ServerKey = viper.GetString("SERVER_KEY")
		midtrans.Environment = midtrans.Sandbox
	}

	return config, err
}
