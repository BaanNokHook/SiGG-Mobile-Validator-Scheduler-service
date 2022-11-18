package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App
		HTTP
		Log
		Mongo
		Redis
		DeviceFilterJobs
	}

	// App -.
	App struct {
		Name    string `env-required:"true" env:"APP_NAME"`
		Version string `env-required:"true" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}

	Mongo struct {
		ConnectionUri        string `env-required:"true" env:"MONGO_URI"`
		Database             string `env-required:"true" env:"MONGO_DB"`
		DeviceCollectionName string `env-required:"true" env:"MONGO_DEVICE_COLLECTION_NAME"`
	}

	//Redis
	Redis struct {
		Addr            string `env-required:"true" env:"REDIS_ADDR"`
		Password        string `env-required:"true" env:"REDIS_PASSWORD"`
		RedisDeviceDB   int    `env-required:"true" env:"REDIS_DEVICE_DB"`
		RedisScheduleDB int    `env-required:"true" env:"REDIS_SCHEDULE_DB"`
	}

	DeviceFilterJobs struct {
		IsRunJob                      bool   `env-required:"true" env:"JOB_RUN"`
		JobName                       string `env-required:"true" env:"JOB_NAME"`
		ScheduleNamespace             string `env-required:"true" env:"JOB_NAMESPACE"`
		ScheduleCron                  string `env-required:"true" env:"JOB_CRON"`
		DeviceOfflineThresholdMinute  int    `env-required:"true" env:"JOB_DEVICE_OFFLINE_THRESHOLD_MINUTE"`
		DeviceValidateThresholdMinute int    `env-required:"true" env:"JOB_DEVICE_VALIDATE_THRESHOLD_MINUTE"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	if _, err := os.Stat(".env"); err == nil {
		err = cleanenv.ReadConfig(".env", cfg)
		if err != nil {
			return nil, fmt.Errorf("config error: %w", err)
		}
	}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
