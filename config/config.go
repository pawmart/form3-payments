package config

import (
	goconfig "github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
)

// DbConfig holding db mongo configuration.
type DbConfig struct {
	Database   string `json:"database,omitempty"`
	Host       string `json:"host,omitempty"`
	User       string `json:"user,omitempty"`
	Password   string `json:"password,omitempty"`
	Auth       string `json:"auth,omitempty"`
	Ssl        string `json:"ssl,omitempty"`
	Replicaset string `json:"replicaset,omitempty"`
}

// Config holding app configuration.
type Config struct {
	Db *DbConfig `protobuf:"bytes,1,opt,name=db,proto3" json:"db,omitempty"`
}

// LoadConfiguration for the app.
func (*Config) LoadConfiguration() *Config {

	// TODO: Fetch config from Vault etc.
	c := new(Config)

	src := env.NewSource(
		env.WithStrippedPrefix("FORM3"),
	)

	conf := goconfig.NewConfig()
	conf.Load(src)
	conf.Scan(c)

	return c
}
