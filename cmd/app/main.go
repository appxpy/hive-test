package main

import (
	"log"

	"github.com/appxpy/hive-test/config"
	"github.com/appxpy/hive-test/internal/app"
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
