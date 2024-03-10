package conf

type config struct {
	Server server `mapstructure:"server"`
	PgSQL  pgsql  `mapstructure:"pgsql"`
	Redis  redis  `mapstructure:"redis"`
	Auth   auth   `mapstructure:"auth"`
	Debug  bool   `mapstructure:"debug"`
}

type server struct {
	Port int    `mapstructure:"port"`
	Site string `mapstructure:"site"`
}

type pgsql struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Auth     string `mapstructure:"auth"`
	Database int    `mapstructure:"database"`
}

type auth struct {
	FrontendCallbackPrefix []string `mapstructure:"frontend_callback_prefix"`
	GitHub                 github   `mapstructure:"github"`
}

type github struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}
