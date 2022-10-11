package processor

import (
	"fmt"
	"os"

	"github.com/riverphillips/envoy-control-plane/api/v1alpha/envoy"
	"github.com/riverphillips/envoy-control-plane/internal"
)

func parseYaml(filepath string) (*envoy.EnvoyConfig, error) {
	cfg := new(envoy.EnvoyConfig)

	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("Error reading YAML file: %s\n", err)
	}

	err = internal.ProtosFromYaml(yamlFile, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
