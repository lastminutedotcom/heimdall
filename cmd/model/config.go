package model

type Config struct {
	CronExpression string         `json:"cron_expression"`
	GraphiteConfig GraphiteConfig `json:"graphite_config"`
}

type GraphiteConfig struct {
	Host string
	Port int
}

func DefautConfig() *Config {
	return &Config{
		CronExpression: "*/5 * * * *",
	}
}
