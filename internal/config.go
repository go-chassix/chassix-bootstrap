package internal

import (
	"strings"
	"sync"

	"gopkg.in/apollo.v0"
)

var (
	simpleConfig   *Config
	chassixConfigs map[string]interface{}
	configs        *Configs
)

const (
	KeyAppConfig     = "chassix.app"
	KeyApolloConfig  = "chassix.apollo"
	KeyServerConfig  = "chassix.server"
	KeyLoggingConfig = "chassix.logging"
)

type Configs struct {
	lock           sync.RWMutex `yaml:"-"`
	chassixConfigs map[string]interface{}
}

//Config all config
type Config struct {
	AppConfig `yaml:",inline"`
	//Server ServerConfig `yaml:"server"`
	ApolloConfig `yaml:",inline"`
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
	Data struct {
		Name    string
		Version string
		Env     string
		Debug   bool
	} `yaml:"app"`
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
func Logging() *LoggingConfig {
	configs.Lock()
	defer configs.UnLock()
	return configs.Logging()
}

// ApolloConfig apollo config
type ApolloConfig struct {
	Data struct {
		Enabled  bool        `yaml:"enabled"`
		Settings apollo.Conf `yaml:"settings"`
	} `yaml:"apollo"`
}

//UsingYaml
func UsingYaml() {
	configs.RLock()
	defer configs.RUnlock()
	LoadConfigsFromEnvFile(configs.chassixConfigs)
}

func UsingApollo() {
	configs.RLock()
	defer configs.RUnlock()
}
func LoadSimpleConfig() {
	LoadFromEnvFile(simpleConfig)
}

// IsApolloEnabled is apollo enable
func IsApolloEnabled() bool {
	return simpleConfig.ApolloConfig.Data.Enabled
}
func Load() {
	LoadSimpleConfig()
	if simpleConfig.AppConfig.Data.Debug {
		IsDebug = true
	}
	if IsApolloEnabled() {
		GetLogger().Info("using apollo config")
		UsingApollo()
		return
	}
	UsingYaml()
	GetLogger().Info("using yaml file config")
}

func init() {
	configsInstance()
	appConfig = new(AppConfig)
	serverConfig = new(ServerConfig)
	apolloConfig = new(ApolloConfig)
	simpleConfig = new(Config)
	loggingConfig = new(LoggingConfig)
}

var appConfig *AppConfig
var serverConfig *ServerConfig
var apolloConfig *ApolloConfig
var loggingConfig *LoggingConfig

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
	configs.Watch(KeyLoggingConfig, loggingConfig)

}

func Watch(name string, cfg interface{}) {
	configs.Watch(name, cfg)
}

func (c *Configs) Watch(name string, cfg interface{}) {
	c.chassixConfigs[name] = cfg
}

//App return app config
func (c *Configs) App() *AppConfig {
	return c.chassixConfigs[KeyAppConfig].(*AppConfig)
}

//Server get server config: port addr
func (c *Configs) Server() *ServerConfig {
	return c.chassixConfigs[KeyServerConfig].(*ServerConfig)
}

//Apollo apollo server settings: ip, namespaces...
func (c *Configs) Apollo() *ApolloConfig {
	return c.chassixConfigs[KeyApolloConfig].(*ApolloConfig)
}

//Get get config by map key
func (c *Configs) Get(key string) (interface{}, bool) {
	val, ok := c.chassixConfigs[key]
	return val, ok
}

//Logging logging config
func (c *Configs) Logging() *LoggingConfig {
	return c.chassixConfigs[KeyLoggingConfig].(*LoggingConfig)
}

type LoggingConfig struct {
	Data struct {
		Level        string
		ReportCaller bool `yaml:"report_caller"`
		NoColors     bool `yaml:"no_colors"`
		CallerFirst  bool `yaml:"caller_first"`
	} `yaml:"logging"`
}

func (logCfg *LoggingConfig) Level() int {
	levelStr := strings.ToLower(logCfg.Data.Level)
	switch levelStr {
	case "debug":
		return 5
	case "info":
		return 4
	case "warn":
		return 3
	case "error":
		return 2
	default:
		return 4
	}
}

func (logCfg *LoggingConfig) ReportCaller() bool {
	return logCfg.Data.ReportCaller
}

func (logCfg *LoggingConfig) NoColors() bool {
	return logCfg.Data.NoColors
}

func (logCfg *LoggingConfig) CallerFirst() bool {
	return logCfg.Data.CallerFirst
}
