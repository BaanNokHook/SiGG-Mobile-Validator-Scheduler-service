package main

import (
	"log"
	"nextclan/validator-register/mobile-validator-scheduler-service/config"
	"nextclan/validator-register/mobile-validator-scheduler-service/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
