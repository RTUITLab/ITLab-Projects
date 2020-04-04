package main

import (
	"ITLab-Projects/config"
	"ITLab-Projects/server"
	"fmt"
)

func main() {
	cfg := config.GetConfig()
	app := &server.App{}
	app.Init(cfg)
	app.Run(":"+cfg.App.AppPort)
	fmt.Scanln()
}
