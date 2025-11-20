package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Romm struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
	} `yaml:"romm"`
	Theme struct {
		FontPath string `yaml:"font_path"`
	} `yaml:"theme"`
}

func New(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
