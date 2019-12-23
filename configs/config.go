package config

import (
	"github.com/BurntSushi/toml"
)

type ConfigTree struct {
	Database DatabaseConfig `toml:"database"`
	Server   ServerConfig   `toml:"server"`
}

type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	SSL      string `toml:"ssl"`
}

type ServerConfig struct {
	Port int `toml:"port"`
}

func LoadConfig(configpath string) (*ConfigTree, error) {
	var (
		Config *ConfigTree
		err    error
	)
	Config = new(ConfigTree)
	_, err = toml.DecodeFile(configpath, Config)
	if err != nil {
		return nil, err
	}
	return Config, nil
}
