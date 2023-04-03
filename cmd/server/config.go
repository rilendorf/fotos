package main

import (
	"gopkg.in/yaml.v3"

	_ "embed"
	"log"
	"os"
	"path"
)

//go:embed default.cfg
var defaultConfig string

type Config struct {
	Listen string `yaml:"Listen"`

	PublicAccess string `yaml:"PublicAccess"`
	DBPath       string `yaml:"DBPath"`
}

func ReadConfig() Config {
	var file *os.File

	file, err := os.Open(*ConfigPath)
	if err != nil {
		log.Printf("Failed to open configuration file '%s': %s; Writing default config", *ConfigPath, err)

		file = DefaultConfig()
	}

	defer file.Close()

	dec := yaml.NewDecoder(file)

	var config *Config
	err = dec.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode yaml configuration from '%s': %s", *ConfigPath, err)
	}

	return *config
}

func DefaultConfig() *os.File {
	path, _ := path.Split(*ConfigPath)
	err := os.MkdirAll(path, 0655)
	if err != nil {
		log.Fatalf("Failed to create configuration Folder at '%s': %s", path, err)
	}

	file, err := os.OpenFile(*ConfigPath, os.O_RDWR|os.O_CREATE, 0655)
	if err != nil {
		log.Fatalf("Failed to open configuration file (readwrite) '%s': %s", *ConfigPath, err)
	}

	_, err = file.WriteString(defaultConfig)
	if err != nil {
		log.Fatalf("Failed to write default configuration '%s': %s", *ConfigPath, err)
	}

	file.Seek(0, 0)

	return file
}
