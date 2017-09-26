package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	. "bosh-generic-sb/broker"
)

type Config struct {
	HTTP     HTTPConfig
	Director DirectorConfig
	Broker   BrokerConfig
}

type HTTPConfig struct {
	Username string
	Password string

	Host string
	Port string
}

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

	err = cfg.Validate()
	if err != nil {
		return cfg, fmt.Errorf("validating config: %s", err)
	}

	return cfg, nil
}

func (c Config) Validate() error {
	err := c.HTTP.Validate()
	if err != nil {
		return fmt.Errorf("validating HTTP: %s", err)
	}

	err = c.Director.Validate()
	if err != nil {
		return fmt.Errorf("validating Director: %s", err)
	}

	err = c.Broker.Validate()
	if err != nil {
		return fmt.Errorf("validating Broker: %s", err)
	}

	return nil
}

func (c HTTPConfig) Validate() error {
	if len(c.Username) == 0 {
		return fmt.Errorf("expected non-empty 'Username'")
	}
	if len(c.Password) == 0 {
		return fmt.Errorf("expected non-empty 'Password'")
	}
	if len(c.Host) == 0 {
		return fmt.Errorf("expected non-empty 'Host'")
	}
	if len(c.Port) == 0 {
		return fmt.Errorf("expected non-empty 'Port'")
	}
	return nil
}
