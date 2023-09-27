package config

type MysqlConfig struct {
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type ServerConfig struct {
	Name         string `mapstructure:"name"`
	MysqlConfig  `mapstructure:"mysql"`
	ConsulConfig `mapstructure:"consul"`
}

var (
	TheServerConfig ServerConfig
)
