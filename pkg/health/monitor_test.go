package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"dock_pilot/internal/config"
	"dock_pilot/pkg/services"
)

func TestCheckServiceHealthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		if _, err := w.Write([]byte("ok")); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	port, _ := strconv.Atoi(u.Port())

	svc := config.ServiceConfig{
		Name:        "svc",
		Port:        port,
		Healthcheck: u.Path,
	}
	mon := NewMonitor(&services.DockerService{}, &config.Config{Services: []config.ServiceConfig{svc}}, nil)
	res := mon.CheckService(context.Background(), svc)
	if res.Status != Healthy {
		t.Errorf("expected healthy, got %s", res.Status)
	}
}
