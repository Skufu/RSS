package main

import (
	"fmt"
	"log"

	"github.com/Skufu/RSS/internal/config"
)

func main() {
	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	fmt.Printf("Initial config: %+v\n", cfg)
	fmt.Printf("Database URL: %s\n", cfg.DatabaseURL)

	// Set the current user and update the config file
	err = cfg.SetUser("Adrian")
	if err != nil {
		log.Fatalf("Error setting user: %v", err)
	}
	fmt.Println("Updated user to 'Adrian'")

	// Read the config file again and print contents
	updatedCfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading updated config: %v", err)
	}
	fmt.Printf("Updated config: %+v\n", updatedCfg)
	fmt.Printf("Database URL: %s\n", updatedCfg.DatabaseURL)
}
