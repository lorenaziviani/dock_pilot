package main

import (
	"context"
	"dock_pilot/internal/config"
	"dock_pilot/pkg/health"
	"dock_pilot/pkg/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: orchestrator <start|stop|restart|status|monitor> [service|all]")
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
	case "monitor":
		logger := log.New(os.Stdout, "", log.LstdFlags)
		monitor := health.NewMonitor(dockerSvc, cfg, logger)
		fmt.Println("[INFO] Starting health monitoring loop...")
		monitor.MonitorLoop(ctx, 10*time.Second)
	case "dashboard":
		fmt.Println("[INFO] Starting DockPilot dashboard (TUI)...")
		RunTUI(ctx, cfg, dockerSvc)
	case "dump":
		fmt.Println("[INFO] Exporting config and service state...")
		file, err := os.Create("dockpilot_dump.json")
		if err != nil {
			log.Fatalf("failed to create dump: %v", err)
		}
		defer file.Close()
		statusMap := make(map[string]string)
		for _, svc := range cfg.Services {
			status, _ := dockerSvc.ContainerStatus(ctx, svc.Name)
			statusMap[svc.Name] = status
		}
		type Dump struct {
			Config interface{}       `json:"config"`
			Status map[string]string `json:"status"`
			Time   string            `json:"time"`
		}
		dump := Dump{Config: cfg, Status: statusMap, Time: time.Now().Format(time.RFC3339)}
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		if err := enc.Encode(dump); err != nil {
			log.Fatalf("failed to export dump: %v", err)
		}
		fmt.Println("[OK] Dump saved to dockpilot_dump.json")
	case "metrics":
		fmt.Println("[INFO] Exposing Prometheus metrics at http://localhost:2112/metrics ...")
		http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			for _, svc := range cfg.Services {
				status, _ := dockerSvc.ContainerStatus(ctx, svc.Name)
				fmt.Fprintf(w, "dockpilot_service_status{service=\"%s\"} %d\n", svc.Name, statusToMetric(status))
			}
		})
		srv := &http.Server{
			Addr:              ":2112",
			Handler:           nil,
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			IdleTimeout:       60 * time.Second,
		}
		log.Fatal(srv.ListenAndServe())
	default:
		fmt.Println("Command not recognized. Use: start, stop, restart, status, monitor, dashboard")
	}
}

func statusToMetric(status string) int {
	switch status {
	case "running":
		return 1
	case "exited":
		return 0
	default:
		return -1
	}
}
