package conf

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/mitchellh/mapstructure"
)

var C config

func Init() error {
	c := viper.New()

	// SET DEFAULT VALUES
	c.SetDefault("auth.bcrypt_rounds", 10)
	c.SetDefault("server.port", 3000)
	c.SetDefault("limit.media.file_size", 51200)
	c.SetDefault("limit.media.image_size", 7680)

	c.SetEnvPrefix("doulog")
	c.SetConfigName("config")
	c.SetConfigType("toml")
	c.AddConfigPath(".")
	c.AddConfigPath("etc")
	c.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
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

	// print bcrypt rounds
	println("Bcrypt rounds: ", C.Auth.BcryptRounds)

	// print all Frontend callback prefix
	for _, prefix := range C.Auth.FrontendCallbackPrefix {
		println("Callback prefix: ", prefix)
	}

	// print postgres
	println("Postgres host: ", C.PgSQL.Host)
	println("Postgres port: ", C.PgSQL.Port)
	println("Postgres username: ", C.PgSQL.Username)
	println("Postgres password: ", C.PgSQL.Password)
	println("Postgres database: ", C.PgSQL.Database)

	// print redis
	println("Redis host: ", C.Redis.Host)
	println("Redis port: ", C.Redis.Port)
	println("Redis auth: ", C.Redis.Auth)
	println("Redis database: ", C.Redis.Database)

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
