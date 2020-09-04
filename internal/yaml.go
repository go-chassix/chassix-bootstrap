package internal

import (
	"io/ioutil"
	"os"
	"sync"

	"c5x.io/logx"
	"gopkg.in/yaml.v3"
)

var ConfigFileEnvKey = "CHASSIX_CONF"

var IsDebug bool
var log *logx.Entry
var logInitOnce sync.Once

func GetLogger() *logx.Entry {
	logInitOnce.Do(func() {
		logger := logx.New()

		if IsDebug {
			logger.ReportCaller = true
			logger.Level = 5
		}
		log = logger.Component("config").Category("boot")
	})
	return log
}

//LoadEnvFile Load config from the file that path is saved in os env.
func LoadFromEnvFile(cfg interface{}) {
	fileName := os.Getenv(ConfigFileEnvKey)
	if err := LoadFromFile(cfg, fileName); err != nil {
		GetLogger().Errorf("load file config error: %s", err)
		os.Exit(1)
	}
}

func LoadChassixConfigs(cfgs []interface{}) {
	fileName := os.Getenv(ConfigFileEnvKey)
	for _, cfg := range cfgs {
		if err := LoadFromFile(cfg, fileName); err != nil {
			GetLogger().Errorf("load file config error: %s\n", err)
			os.Exit(1)
		}
	}

}

//LoadConfigsFromFile from file
func LoadConfigsFromEnvFile(cfgs map[string]interface{}) error {
	fileName := os.Getenv(ConfigFileEnvKey)
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	data, _ := ioutil.ReadAll(f)
	GetLogger().Debugf("==========config file=========\n%s", data)
	for key, cfgPtr := range cfgs {
		if err := yaml.Unmarshal(data, cfgPtr); err != nil {
			GetLogger().Error("error", err.Error())
		}

		GetLogger().Debugf("load key: %s,value: %+v\n", key, cfgPtr)
	}
	return nil
}

//LoadFromFile from file
func LoadFromFile(cfg interface{}, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return err
	}
	return nil
}

//LoadCustomFromFile Load custom config from file, save to custom config
func LoadCustomFromFile(customCfg interface{}, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(customCfg); err != nil {
		return err
	}
	return nil
}
