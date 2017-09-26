package broker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	HTTP HTTPConfig

	ServiceInstanceManifest string
	ServiceInstanceParams   []Param
	ServiceBindingManifest  string
	ServiceBindingParams    []Param

	Director DirectorConfig
}

type HTTPConfig struct {
	Username string
	Password string

	Host string
	Port string
}

type Param struct {
	Name string
	Ops  []interface{}
}

func (c Config) AllowsServiceBinding() bool { return len(c.ServiceBindingManifest) > 0 }

func NewConfigFromPath(path string) (Config, error) {
	var cfg Config

	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("reading job config: %s", err)
	}

	err = json.Unmarshal(configBytes, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("unmarshaling job config: %s", err)
	}

	return cfg, nil
}
