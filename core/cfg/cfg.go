package cfgcore

import (
	"github.com/caarlos0/env/v5"
	"github.com/joho/godotenv"
)

const (
	MemCacheKeyReloadConfig = "cacheKey:reload-config"
)

// Load loads configuration from local .env file
func LoadLocal(out interface{}, stage string) error {
	if err := PreloadLocalENV(stage); err != nil {
		return err
	}

	if err := env.Parse(out); err != nil {
		return err
	}

	return nil
}

// Load loads configuration from local .env file
func Load(out interface{}, stage string) error {
	if err := PreloadENV(stage); err != nil {
		return err
	}

	if err := env.Parse(out); err != nil {
		return err
	}

	return nil
}

// LoadWithAPS loads configuration from local .env file and AWS Parameter Store as well
func LoadWithAPS(out interface{}, appName, stage string) error {
	if appName != "" && stage != "development" {
		return Load(out, stage)
	}
	return LoadLocal(out, stage)
}

// PreloadLocalENV reads .env* files and sets the values to os ENV
func PreloadLocalENV(stage string) error {
	return godotenv.Load(".env.local")
}

func PreloadENV(stage string) error {
	return godotenv.Load(".env")
}
