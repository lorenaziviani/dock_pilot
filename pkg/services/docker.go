package services

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DockerService struct {
	cli *client.Client
}

func NewDockerService() (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerService{cli: cli}, nil
}

func (d *DockerService) StartContainer(ctx context.Context, name, image string, ports []string, volumes map[string]string) error {
	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}
	for _, p := range ports {
		parts := strings.Split(p, ":")
		if len(parts) == 2 {
			hostPort := parts[0]
			containerPort := parts[1] + "/tcp"
			portKey := nat.Port(containerPort)
			portBindings[portKey] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: hostPort}}
			exposedPorts[portKey] = struct{}{}
		}
	}

	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image:        image,
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		PortBindings: portBindings,
	}, nil, nil, name)
	if err != nil {
		return err
	}
	return d.cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
}

func (d *DockerService) StopContainer(ctx context.Context, name string) error {
	return d.cli.ContainerStop(ctx, name, container.StopOptions{})
}

func (d *DockerService) RestartContainer(ctx context.Context, name string) error {
	return d.cli.ContainerRestart(ctx, name, container.StopOptions{})
}

func (d *DockerService) ContainerStatus(ctx context.Context, name string) (string, error) {
	containerJSON, err := d.cli.ContainerInspect(ctx, name)
	if err != nil {
		return "", err
	}
	return containerJSON.State.Status, nil
}

func (d *DockerService) ListContainers(ctx context.Context) ([]container.Summary, error) {
	return d.cli.ContainerList(ctx, container.ListOptions{All: true})
}
