package model

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type KubeConfig struct {
	MgmtPort string `json:"management_port"`
}

type Config struct {
	GraphiteConfig      GraphiteConfig `json:"graphite_config"`
	CollectEveryMinutes string         `json:"collect_every_minutes"`
	KubeConfig          *KubeConfig    `json:"kubernetes,omitempty"`
}

type GraphiteConfig struct {
	Host string
	Port int
}

func ParseConfig(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var config *Config
	if err := json.Unmarshal([]byte(b), &config); err != nil {
		return nil, err
	}
	return config, nil
}
