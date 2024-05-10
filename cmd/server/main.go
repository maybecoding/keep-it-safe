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

	a := app.New(cfg)
	err = a.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	a.Run()
}
