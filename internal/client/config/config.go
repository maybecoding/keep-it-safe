package config

import (
	"fmt"
	"reflect"

	"github.com/maybecoding/keep-it-safe/pkg/vscfg"
)

type Config struct {
	Server
	TUI
}

type (
	Server struct {
		Address string `default:"http://localhost:8080" flg:"s" flgU:"Endpoint HTTP-server address" env:"SERVER_ADDRESS"`
	}
	TUI struct {
		WindowHeight int `default:"25" flg:"wh" flgU:"TUI interface window height" env:"WINDOW_HEIGHT"`
	}
)

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
