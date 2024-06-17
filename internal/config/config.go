package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadAllConfigs set various configs
func LoadAll(envFile string) {
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("can't load environment file. error: %v", err)
	}

	LoadApp()
	LoadDBCfg()
}
