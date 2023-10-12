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

type MysqlConfig struct {
	UserName string `yaml:"username" json:"username"  mapstructure:"username"`
	Password string `yaml:"password" json:"password"  mapstructure:"password"`
	Host     string `yaml:"host" json:"host"  mapstructure:"host"`
	Port     int    `yaml:"port" json:"port"  mapstructure:"port"`
	DBName   string `yaml:"dbname" json:"dbname"  mapstructure:"dbname"`
}

type ConsulConfig struct {
	Host string `yaml:"host" json:"host"  mapstructure:"host"`
	Port int    `yaml:"port" json:"port"  mapstructure:"port"`
}

type ServerConfig struct {
	Name         string   `yaml:"name" json:"name"  mapstructure:"name"`
	Host         string   `yaml:"host" json:"host"  mapstructure:"host"`
	Tags         []string `yaml:"tags" json:"tags"  mapstructure:"tags"`
	MysqlConfig  `yaml:"mysql" json:"mysql"  mapstructure:"mysql"`
	ConsulConfig `yaml:"consul" json:"consul"  mapstructure:"consul"`
}

var (
	TheServerConfig ServerConfig
	TheNacosConfig  NacosConfig
)
