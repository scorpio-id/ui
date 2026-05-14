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

	return cfg
}
