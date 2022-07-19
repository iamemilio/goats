package generator

import (
	"errors"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type RunConfig struct {
	Version float32 `yaml:"version"`
	Runs    []Run   `yaml:"runs"`
}

type Run struct {
	Name          string `yaml:"name"`
	Kubeconfig    string `yaml:"kubeconfig"`
	App           `yaml:"app"`
	TrafficDriver `yaml:"traffic-driver"`
}

type App struct {
	Image   string            `yaml:"image"`
	EnvVars map[string]string `yaml:"environment-variables"`
}

type TrafficDriver struct {
	Endpoint string `yaml:"service-endpoint"`
	Delay    uint   `yaml:"startup-delay"`
	Traffic  `yaml:"traffic"`
}

type Traffic struct {
	Duration uint `yaml:"duration"`
	Rate     uint `yaml:"requests-per-second"`
	Users    uint `yaml:"concurrent-requests"`
}

func (config *RunConfig) Parse(file []byte) error {
	err := yaml.Unmarshal(file, config)
	if err != nil {
		return err
	}

	if len(config.Runs) == 0 {
		return errors.New("run config must have at least one run")
	}

	for _, run := range config.Runs {
		if run.Kubeconfig == "" {
			run.Kubeconfig, err = filepath.Abs("~/.kube/config")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
