package envs

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	JiraApiToken  string `env:"JIRA_API_TOKEN,required"`
	JiraBaseUrl   string `env:"JIRA_BASE_URL,required"`
	JiraUserEmail string `env:"JIRA_USER_EMAIL,required"`
}

func New() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return cfg, err
}
