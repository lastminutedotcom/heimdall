package model

type Config struct {
	GraphiteConfig      GraphiteConfig `json:"graphite_config"`
	CollectEveryMinutes string         `json:"collect_every_minutes"`
}

type GraphiteConfig struct {
	Host string
	Port int
}

func DefautConfig() *Config {
	return &Config{
		CollectEveryMinutes: "5",
	}
}
