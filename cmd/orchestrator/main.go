package main

import (
	"context"
	"dock_pilot/internal/config"
	"dock_pilot/pkg/services"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: orchestrator <start|stop|restart|status> [service|all]")
		os.Exit(1)
	}
	cmd := os.Args[1]
	arg := "all"
	if len(os.Args) > 2 {
		arg = os.Args[2]
	}

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	dockerSvc, err := services.NewDockerService()
	if err != nil {
		log.Fatalf("failed to connect to Docker: %v", err)
	}
	ctx := context.Background()

	switch cmd {
	case "start":
		for _, svc := range cfg.Services {
			if arg == "all" || arg == svc.Name {
				err := dockerSvc.StartContainer(ctx, svc.Name, svc.Image, nil, nil)
				if err != nil {
					fmt.Printf("[ERROR] %s: %v\n", svc.Name, err)
				} else {
					fmt.Printf("[OK] Started %s\n", svc.Name)
				}
			}
		}
	case "stop":
		for _, svc := range cfg.Services {
			if arg == "all" || arg == svc.Name {
				err := dockerSvc.StopContainer(ctx, svc.Name)
				if err != nil {
					fmt.Printf("[ERROR] %s: %v\n", svc.Name, err)
				} else {
					fmt.Printf("[OK] Stopped %s\n", svc.Name)
				}
			}
		}
	case "restart":
		for _, svc := range cfg.Services {
			if arg == "all" || arg == svc.Name {
				err := dockerSvc.RestartContainer(ctx, svc.Name)
				if err != nil {
					fmt.Printf("[ERROR] %s: %v\n", svc.Name, err)
				} else {
					fmt.Printf("[OK] Restarted %s\n", svc.Name)
				}
			}
		}
	case "status":
		for _, svc := range cfg.Services {
			if arg == "all" || arg == svc.Name {
				status, err := dockerSvc.ContainerStatus(ctx, svc.Name)
				if err != nil {
					fmt.Printf("[ERROR] %s: %v\n", svc.Name, err)
				} else {
					fmt.Printf("%s: %s\n", svc.Name, status)
				}
			}
		}
	default:
		fmt.Println("Command not recognized. Use: start, stop, restart, status")
	}
}
