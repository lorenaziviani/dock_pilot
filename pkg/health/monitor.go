package health

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"dock_pilot/internal/config"
	"dock_pilot/pkg/services"
)

type HealthStatus string

const (
	Healthy     HealthStatus = "healthy"
	Degraded    HealthStatus = "degraded"
	Unreachable HealthStatus = "unreachable"
)

type ServiceHealth struct {
	Name   string
	Status HealthStatus
	Detail string
}

type Monitor struct {
	Docker *services.DockerService
	Config *config.Config
	Logger *log.Logger
}

func NewMonitor(docker *services.DockerService, cfg *config.Config, logger *log.Logger) *Monitor {
	return &Monitor{Docker: docker, Config: cfg, Logger: logger}
}

func (m *Monitor) CheckService(ctx context.Context, svc config.ServiceConfig) ServiceHealth {
	url := fmt.Sprintf("http://localhost:%d%s", svc.Port, svc.Healthcheck)
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return ServiceHealth{svc.Name, Unreachable, err.Error()}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		return ServiceHealth{svc.Name, Healthy, string(body)}
	}
	return ServiceHealth{svc.Name, Degraded, fmt.Sprintf("Status: %d, Body: %s", resp.StatusCode, string(body))}
}

func (m *Monitor) MonitorLoop(ctx context.Context, interval time.Duration) {
	for {
		for _, svc := range m.Config.Services {
			logger, err := services.NewServiceLogger(svc.Name, "./logs")
			if err != nil {
				m.Logger.Printf("[ERROR] Falha ao criar logger para %s: %v", svc.Name, err)
				continue
			}

			health := m.CheckService(ctx, svc)
			logger.Log(fmt.Sprintf("[%s] %s - %s", health.Status, svc.Name, health.Detail))
			m.Logger.Printf("[%s] %s - %s", health.Status, svc.Name, health.Detail)

			if health.Status == Unreachable {
				logger.Log(fmt.Sprintf("[ACTION] Restarting %s...", svc.Name))
				m.Logger.Printf("[ACTION] Restarting %s...", svc.Name)
				err := m.Docker.RestartContainer(ctx, svc.Name)
				if err != nil {
					logger.Log(fmt.Sprintf("[ERROR] Failed to restart %s: %v", svc.Name, err))
					m.Logger.Printf("[ERROR] Failed to restart %s: %v", svc.Name, err)
				}
			}
			logger.Close()
		}
		time.Sleep(interval)
	}
}
