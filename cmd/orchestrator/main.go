package main

import (
	"dock_pilot/internal/config"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	fmt.Println("Loaded services:")
	for _, svc := range cfg.Services {
		fmt.Printf("- %s (%s) on port %d, health: %s\n", svc.Name, svc.Image, svc.Port, svc.Healthcheck)
	}
}
