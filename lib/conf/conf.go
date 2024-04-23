package conf

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/mitchellh/mapstructure"
)

var C config

func Init() error {
	c := viper.New()
	c.SetEnvPrefix("doulog")
	c.SetConfigName("config")
	c.SetConfigType("toml")
	c.AddConfigPath(".")
	c.AddConfigPath("etc")
	replacer := strings.NewReplacer(".", "_")
	c.SetEnvKeyReplacer(replacer)
	c.AutomaticEnv()

	// ADD START
	envKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(C, &envKeysMap); err != nil {
		return err
	}
	for k := range *envKeysMap {
		if bindErr := viper.BindEnv(k); bindErr != nil {
			return bindErr
		}
	}
	// ADD END

	// If a config file is found, read it in.
	if err := c.ReadInConfig(); err != nil {
		return err
	}

	if err := c.Unmarshal(&C); err != nil {
		return err
	}
	return nil
}
