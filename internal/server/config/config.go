// Package config for configuration of server side application.
package config

import (
	"fmt"
	"reflect"
	"time"

	"github.com/maybecoding/keep-it-safe/pkg/vscfg"
)

// Config type for configuration.
type Config struct {
	JWT
	DB
	HTTP
	Encryption
	Log
}

type (
	// JWT struct for JWT settings.
	JWT struct {
		Secret       string `default:"super complex secret nobody can read it" flg:"jwtsecret" flgU:"jwt secret" env:"JWT_SECRET"`
		ExpiresHours int    `default:"24"`
	}
	// DB - struct for db config.
	DB struct { // postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable
		Path string `default:"postgres://api:pwd@localhost:5432/keep_it_safe?sslmode=disable" flg:"dbpath" flgU:"postgres database connection" env:"DB_PATH"`
	}
	// HTTP - struct for HTTP settings.
	HTTP struct {
		Address string `default:"localhost:8080" flg:"http" flgU:"Endpoint HTTP-server address" env:"HTTP_ADDRESS"`
	}
	// Encryption - struct for encryption settings.
	Encryption struct {
		MasterKeyHex   string        `default:"AAAAAAAAAABBBBBBBBBBCCCCCCCCCCAAAAAAAAAABBBBBBBBBBCCCCCCCCCCAAAA" flg:"encr_key" env:"ENCR_KEY"`
		RotateDuration time.Duration `default:"10m"`
		KeySize        int           `default:"32"`
	}
	// Log - logger configuration.
	Log struct {
		Level string `default:"" flg:"log" flgU:"Log level" env:"LOG_LEVEL"`
	}
)

// New - creates new config and fiiles it based on defaults, env and flags.
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
