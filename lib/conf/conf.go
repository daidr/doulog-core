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
	keys := getAllKeys(C, "")
	for _, key := range keys {
		err := c.BindEnv(key[:len(key)-1])
		if err != nil {
			return err
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

func getAllKeys(C interface{}, parentKeyChain string) []string {
	var keys []string
	if parentKeyChain != "" {
		parentKeyChain += "."
	}
	// 如果是数组，直接返回
	if _, ok := C.([]string); ok {
		keys = append(keys, parentKeyChain)
		return keys
	}
	tempKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(C, &tempKeysMap); err != nil {
		keys = append(keys, parentKeyChain)
	} else {
		for key := range *tempKeysMap {
			keys = append(keys, getAllKeys((*tempKeysMap)[key], parentKeyChain+key)...)
		}
	}
	return keys
}
