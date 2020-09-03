package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var ConfigFileEnvKey = "CHASSIX_CONF"

const (
	SourceApollo = "apollo"
	SourceYaml   = "yaml"
	SourceYml    = "yml"
)

//LoadEnvFile Load config from the file that path is saved in os env.
func LoadFromEnvFile(cfg interface{}) {
	fileName := os.Getenv(ConfigFileEnvKey)
	if err := LoadFromFile(cfg, fileName); err != nil {
		fmt.Printf("load file config error: %s\n", err)
		os.Exit(1)
	}
}

func LoadChassixConfigs(cfgs []interface{}) {
	fileName := os.Getenv(ConfigFileEnvKey)
	for _, cfg := range cfgs {
		if err := LoadFromFile(cfg, fileName); err != nil {
			fmt.Printf("load file config error: %s\n", err)
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
	fmt.Println(string(data))
	for key, cfgPtr := range cfgs {

		//val:= reflect.New(reflect.ValueOf(cfgs["chassix.server"]).Elem().Type())
		if err := yaml.Unmarshal(data, cfgPtr); err != nil {
			fmt.Println("error", err.Error())
		}
		fmt.Printf("load key: %s,value: %+v\n", key, cfgPtr)
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
