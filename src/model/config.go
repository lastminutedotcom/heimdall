package model

type Config struct {
	CronExpression string `json:"cron_expression"`
}

func DefautConfig() *Config {
	return &Config{
		CronExpression: "*/5 * * * *",
	}
}
