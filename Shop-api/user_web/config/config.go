package config

// UserSrvConfig  user grpc配置
type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

// UserRedisConfig user redis 配置
type UserRedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

// ServerConfig  总配置
type ServerConfig struct {
	Name             string          `mapstructure:"name"`
	UserSrvInfo      UserSrvConfig   `mapstructure:"user-srv"`
	UserRedisInfo    UserRedisConfig `mapstructure:"redis"`
	ConsulConfigInfo ConsulConfig    `mapstructure:"consul"`
}

var (
	//TheServerConfig 总配置文件
	TheServerConfig ServerConfig
)
