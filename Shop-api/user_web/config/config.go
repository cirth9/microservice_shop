package config

type NacosConfig struct {
	NacosServer NacosServer `yaml:"nacos_server" mapstructure:"nacos_server"`
	NacosClient NacosClient `yaml:"nacos_client" mapstructure:"nacos_client"`
}

type NacosServer struct {
	DataId string `yaml:"dataId" mapstructure:"dataId"`
	Ip     string `yaml:"ip"mapstructure:"ip"`
	Port   uint64 `yaml:"port"mapstructure:"port"`
}

type NacosClient struct {
	NotLoadCacheAtStart bool   `yaml:"not_load_cache_at_start" mapstructure:"not_load_cache_at_start"`
	LogDir              string `yaml:"log_dir" mapstructure:"log_dir"`
	CacheDir            string `yaml:"cache_dir" mapstructure:"cache_dir"`
	NamespaceId         string `yaml:"namespace_id" mapstructure:"namespace_id"`
	TimeoutMs           uint64 `yaml:"timeout_ms" mapstructure:"timeout_ms"`
}

// UserSrvConfig  user grpc配置
type UserSrvConfig struct {
	//Host string `yaml:"host" mapstructure:"host"`
	//Port int    `yaml:"port" mapstructure:"port"`
	Name string   `yaml:"name" mapstructure:"name" json:"name"`
	Tags []string `yaml:"tags" mapstructure:"tags" json:"tags"`
}

// UserRedisConfig user redis 配置相关信息
type UserRedisConfig struct {
	Host string `yaml:"host" mapstructure:"host" json:"host"`
	Port int    `yaml:"port" mapstructure:"port" json:"port"`
}

// ConsulConfig 主要用于链接consul的相关配置
type ConsulConfig struct {
	Host string `yaml:"host" mapstructure:"host" json:"host"`
	Port int    `yaml:"port" mapstructure:"port" json:"port"`
	//Name string `yaml:"name" mapstructure:"name"`
}

// ServerConfig  总配置
type ServerConfig struct {
	Name             string          `yaml:"name" mapstructure:"name" json:"host"`
	UserSrvInfo      UserSrvConfig   `yaml:"user-srv" mapstructure:"user-srv" json:"user-srv"`
	UserRedisInfo    UserRedisConfig `yaml:"redis" mapstructure:"redis" json:"redis"`
	ConsulConfigInfo ConsulConfig    `yaml:"consul" mapstructure:"consul" json:"consul"`
}

var (
	//TheServerConfig 总配置文件
	TheServerConfig ServerConfig
	TheNacosConfig  NacosConfig
)
