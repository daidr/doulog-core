package conf

import (
	"github.com/spf13/viper"
)

var C config

func Init() error {
	c := viper.New()
	c.SetConfigName("config")
	c.AddConfigPath(".")
	c.AddConfigPath("etc")

	// If a config file is found, read it in.
	if err := c.ReadInConfig(); err != nil {
		return err
	}

	if err := c.Unmarshal(&C); err != nil {
		return err
	}
	return nil
}
