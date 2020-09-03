package config

import (
	"sync"

	"c5x.io/bootstrap/config/internal"
	"gopkg.in/apollo.v0"
)

var (
	simpleConfig   *Config
	chassixConfigs map[string]interface{}
	configs        *Configs
)

const (
	KeyAppConfig        = "chassix.app"
	KeyApolloConfig     = "chassix.apollo"
	KeyServerConfig     = "chassix.server"
	KeyDatasourceConfig = "chassix.datasource"
	KeyRedisConfig      = "chassix.redis"
)

type Configs struct {
	lock           sync.RWMutex `yaml:"-"`
	chassixConfigs map[string]interface{}
}

//Config all config
type Config struct {
	App AppConfig `yaml:"app"`
	//Server ServerConfig `yaml:"server"`
	Apollo ApolloConfig `yaml:"apollo"`
	//lock   sync.RWMutex `yaml:"-"`
}

func (c *Configs) RLock() {
	c.lock.RLock()
}
func (c *Configs) RUnlock() {
	c.lock.RUnlock()
}

func (c *Configs) Lock() {
	c.lock.Lock()
}
func (c *Configs) UnLock() {
	c.lock.Unlock()
}

//AppConfig application config
type AppConfig struct {
	Name    string
	Version string
	Env     string
	Config  string
}

//ServerConfig server config
type ServerConfig struct {
	Data ServerConfigData `yaml:"server"`
}

//ServerConfig
type ServerConfigData struct {
	Port int
	Addr string
}

//App app config
func App() *AppConfig {
	configs.Lock()
	defer configs.UnLock()
	return configs.App()
}

//Server server config
func Server() *ServerConfig {
	configs.Lock()
	defer configs.UnLock()
	return configs.Server()
}

// ApolloConfig apollo config
type ApolloConfig struct {
	Enabled  bool        `yaml:"enabled"`
	Settings apollo.Conf `yaml:"settings"`
}

//UsingYaml
func UsingYaml() {
	configs.RLock()
	defer configs.RUnlock()
	internal.LoadConfigsFromEnvFile(configs.chassixConfigs)
}

func UsingApollo() {
	configs.RLock()
	defer configs.RUnlock()
}
func LoadSimpleConfig() {
	internal.LoadFromEnvFile(simpleConfig)
}

// IsApolloEnabled is apollo enable
func IsApolloEnabled() bool {
	return simpleConfig.Apollo.Enabled
}
func Load() {
	LoadSimpleConfig()
	if IsApolloEnabled() {
		UsingApollo()
		return
	}
	UsingYaml()
}

func init() {
	configsInstance()
	appConfig = new(AppConfig)
	serverConfig = new(ServerConfig)
	apolloConfig = new(ApolloConfig)
	simpleConfig = new(Config)
}

var appConfig *AppConfig
var serverConfig *ServerConfig
var apolloConfig *ApolloConfig

var configsInitOnce sync.Once

func configsInstance() *Configs {
	configsInitOnce.Do(func() {
		configs = new(Configs)
		configs.chassixConfigs = make(map[string]interface{})
	})
	return configs
}

func WatchBootstrapConfig() {
	configs.Watch(KeyAppConfig, appConfig)
	configs.Watch(KeyServerConfig, serverConfig)
	configs.Watch(KeyApolloConfig, apolloConfig)

}
func (c *Configs) Watch(name string, cfg interface{}) {
	c.chassixConfigs[name] = cfg
}

func (c *Configs) App() *AppConfig {
	return c.chassixConfigs[KeyAppConfig].(*AppConfig)
}
func (c *Configs) Server() *ServerConfig {
	return c.chassixConfigs[KeyServerConfig].(*ServerConfig)
}
func (c *Configs) Apollo() *ApolloConfig {
	return c.chassixConfigs[KeyApolloConfig].(*ApolloConfig)
}
