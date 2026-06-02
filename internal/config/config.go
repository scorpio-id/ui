package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config provides a template for marshalling .yml configuration files
type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	SPNEGO struct {
		Realm                string `yaml:"realm" json:"realm"`
		ServicePrincipalName string `yaml:"service_principal_name" json:"service_principal_name"`
		Password             string `yaml:"password" json:"-"`
	} `yaml:"spnego" json:"spnego"`
	PKI struct {
		Endpoint             string   `yaml:"endpoint" json:"endpoint"`
		ServicePrincipalName string   `yaml:"service_principal_name" json:"service_principal_name"`
		SANs                 []string `yaml:"sans" json:"sans"`
	} `yaml:"pki" json:"pki"`
	Persistence struct {
		Enabled  bool   `yaml:"enabled" json:"enabled"`
		Port     string `yaml:"port" json:"port"`
		Host     string `yaml:"host" json:"host"`
		User     string `yaml:"user" json:"user"`
		Path     string `yaml:"path" json:"path"`
		Password string `yaml:"-" json:"-"` // DO NOT MARSHAL PASSWORD!
		Database int    `yaml:"database" json:"database"`
	} `yaml:"persistence" json:"persistence"`
}

// NewConfig takes a .yml filename from the same /config directory, and returns a populated configuration
func NewConfig(s string) Config {
	f, err := os.Open(s)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// TODO retrieve content from Kube Secrets using configured file paths if persistence enabled
	if cfg.Persistence.Enabled {
		content, err := os.ReadFile(cfg.Persistence.Path)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		cfg.Persistence.Password = string(content)
	}

	return cfg
}
