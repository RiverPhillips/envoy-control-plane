package processor

import (
	"fmt"
	"os"

	"github.com/ghodss/yaml"

	"github.com/riverphillips/envoy-control-plane/api/v1alpha/envoy"
)

func parseYaml(filepath string) (*envoy.EnvoyConfig, error) {
	var cfg envoy.EnvoyConfig

	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("Error reading YAML file: %s\n", err)
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
