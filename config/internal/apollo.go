package internal

import (
	"encoding/json"
	"fmt"
	"gopkg.in/apollo.v0"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type HotReLoader struct {
	Config interface{}
	Func   func(cfg interface{})
}

func LoadFromApollo(cfg interface{}, settings *apollo.Conf, configs []interface{}) {
	//LoadFromEnvFile()
	//if IsApolloEnable() {

	if err := apollo.StartWithConf(settings); err != nil {
		fmt.Printf("load apollo config error: %s\n", err)
		os.Exit(1)
		return
	}
	readConfig(cfg, settings.Namespaces)

	// hot loading refresh config
	go func() {
		for {
			event := apollo.WatchUpdate()
			changeEvent := <-event
			bytes, _ := json.Marshal(changeEvent)
			fmt.Println("event:", string(bytes))

			for _, cfg := range configs {
				readConfig(cfg, settings.Namespaces)
			}
		}
	}()
	//}
}

//LoadCustomFromFile Load custom config from apollo, save to custom config
func LoadCustomFromApollo(customCfg interface{}, settings *apollo.Conf) error {
	//if !IsApolloEnable() {
	//
	//	return errors.New("apollo is not enabled")
	//}
	readConfig(customCfg, settings.Namespaces)
	return nil
}
func readConfig(cfg interface{}, namespaces []string) {
	for _, namespace := range namespaces {
		ymlTxt := apollo.GetNameSpaceContent(namespace, "")
		if err := yaml.NewDecoder(strings.NewReader(ymlTxt)).Decode(cfg); err != nil {
			fmt.Printf("load apollo config error: %s\n", err)
		}
	}
}
