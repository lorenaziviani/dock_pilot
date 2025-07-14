package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	Name        string   `yaml:"name"`
	Image       string   `yaml:"image"`
	Port        int      `yaml:"port"`
	Healthcheck string   `yaml:"healthcheck"`
	Volumes     []string `yaml:"volumes"`
	Networks    []string `yaml:"networks"`
	Ports       []string `yaml:"ports"`
}

type Config struct {
	Services []ServiceConfig `yaml:"services"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
