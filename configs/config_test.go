package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	var (
		toml         *ConfigTree
		sampleConfig *ConfigTree
		err          error
	)

	sampleConfig = &ConfigTree{
		Server: ServerConfig{
			Port: 3030,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "persona",
			Password: "password",
			Database: "persona",
			SSL:      "disable",
		},
	}

	toml, err = LoadConfig("config.sample.toml")
	if err != nil {
		t.Fatalf("func TestLoadConfig() Failed: var err is not nil. detail: %s", err)
	}
	if toml.Server != sampleConfig.Server {
		t.Fatalf("func TestLoadConfig() Failed: [server] does no match. config.toml is %+v", toml)
	}
	if toml.Database != sampleConfig.Database {
		t.Fatalf("func TestLoadConfig() Failed: [database] does no match. config.toml is %+v", toml)
	}
}

func TestConfigFileNotFound(t *testing.T) {
	var (
		toml *ConfigTree
		err  error
	)
	toml, err = LoadConfig("notfound.toml")
	if toml != nil {
		t.Fatalf("func TestLoadConfigError() Failed: var toml is not nil.")
	}
	if err == nil {
		t.Fatalf("func TestLoadConfigError() Failed: var err is nil.")
	}
}
