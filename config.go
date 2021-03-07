package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Change represents a single change to make to copied files
type Change struct {
	Replace string `yaml:"replace"`
	With    string `yaml:"with"`
	Default string `yaml:"default"`
}

// Config represents configuration to copy and modify files
type Config struct {
	Message         string             `yaml:"message"`
	Path            string             `yaml:"path"`
	Changes         map[string]*Change `yaml:"changes"`
	ModifyFilenames bool               `yaml:"modifyfilenames"`
}

func getConfig() *Config {
	bytes, err := ioutil.ReadFile("./.makecopy.yml")
	if err != nil {
		log.Fatal(err)
	}

	c := Config{}
	if err := yaml.Unmarshal(bytes, &c); err != nil {
		log.Fatal(err)
	}

	// TODO(jan25): Validate config

	return &c
}
