package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Parser interface {
	Parse() (Config, error)
}

type YamlParser struct {
	filepath string
}

func NewParser(path string) Parser {
	return YamlParser{filepath: path}
}

func (p YamlParser) Parse() (Config, error) {
	config := Config{}

	data, err := ioutil.ReadFile(p.filepath)
	if err != nil {
		return config, fmt.Errorf("could not open config file %v", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("could not parse config yaml %v", err)
	}

	return config, nil
}
