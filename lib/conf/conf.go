package conf

import (
	"strings"

	"github.com/spf13/viper"
)

var C config

func Init() error {
	c := viper.New()
	c.SetEnvPrefix("doulog")
	c.SetConfigName("config")
	c.AddConfigPath(".")
	c.AddConfigPath("etc")
	replacer := strings.NewReplacer(".", "_")
	c.SetEnvKeyReplacer(replacer)
	c.AutomaticEnv()

	// If a config file is found, read it in.
	if err := c.ReadInConfig(); err != nil {
		return err
	}

	if err := c.Unmarshal(&C); err != nil {
		return err
	}
	return nil
}
