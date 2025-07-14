package services

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

func (d *DockerService) StartContainer(ctx context.Context, name, image string, ports map[string]string, volumes map[string]string) error {
	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image: image,
	}, nil, nil, nil, name)
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
