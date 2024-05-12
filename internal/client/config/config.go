package config

import (
	"fmt"
	"reflect"

	"github.com/maybecoding/keep-it-safe/pkg/vscfg"
)

// Config - Struct with app configuration.
type Config struct {
	Server
	Log
}

type (
	// Server configuration.
	Server struct {
		Address string `default:"http://localhost:8080" flg:"s" flgU:"Endpoint HTTP-server address" env:"SERVER_ADDRESS"`
	}
	// Log - logger configuration.
	Log struct {
		Level string `default:"debug" flg:"log" flgU:"Log level" env:"LOG_LEVEL"`
	}
)

// New creates new log struct and fill it according to tags using vscfg lib.
func New() (*Config, error) {
	cfg := new(Config)
	rCfg := reflect.ValueOf(cfg).Elem()
	// Заполняем значениями по умолчанию, флагами и env
	var fns []vscfg.Fn
	fns = append(fns, vscfg.Tag("default"))
	fns = append(fns, vscfg.Flag("flg", "flgU")...)
	fns = append(fns, vscfg.Env("env"))
	err := vscfg.FillByTags(rCfg, fns...)
	if err != nil {
		return nil, fmt.Errorf("config - New - vscfg.FillByTags: %w", err)
	}
	return cfg, nil
}
