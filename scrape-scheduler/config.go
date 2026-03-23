package scrape_scheduler

import (
	"log"

	"github.com/joho/godotenv"
)

const (
	EnvironmentLocal = "local"
	EnvironmentProd  = "prod"
	EnvironmentDev   = "dev"
)

type Config struct {
	env map[string]string
}

func NewConfig(environment string) *Config {
	var envMap map[string]string
	var err error

	switch environment {
	case EnvironmentLocal:
		envMap, err = godotenv.Read(".env.local")
		if err != nil {
			log.Fatalf("Error loading .env.local file: %v", err)
		}
	case EnvironmentProd:
		envMap, err = godotenv.Read(".env")
		if err != nil {
			log.Fatalf("Error loading .env.prod file: %v", err)
		}
	case EnvironmentDev:
		envMap, err = godotenv.Read(".env.dev")
		if err != nil {
			log.Fatalf("Error loading .env.dev file: %v", err)
		}
	}

	return &Config{env: envMap}
}

func (c *Config) GetAllEnv() map[string]string {
	envCopy := make(map[string]string, len(c.env))
	for k, v := range c.env {
		envCopy[k] = v
	}
	return envCopy
}

func (c *Config) GetTemporalHostPort() string {
	return c.env["TEMPORAL_HOST_PORT"]
}

func (c *Config) GetSchedulerCron() string {
	return c.env["SCHEDULER_CRON"]
}

func (c *Config) GetScraperGatewayHost() string {
	return c.env["SCRAPER_GATEWAY_HOST"]
}

func (c *Config) GetSchedulerTaskQueue() string {
	return c.env["SCHEDULER_TASK_QUEUE"]
}
