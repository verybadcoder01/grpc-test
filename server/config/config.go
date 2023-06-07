package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	DbPath  string `yaml:"db_path"`
	LogPath string `yaml:"log_path"`
}

func ParseConfig() Config {
	var config Config
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("cant parse config %v", err.Error())
	}
	return config
}
