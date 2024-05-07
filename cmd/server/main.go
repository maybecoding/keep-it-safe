package main

import (
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/server/app"
	"github.com/maybecoding/keep-it-safe/internal/server/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	app := app.New(cfg)
	err = app.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	app.Run()
}
