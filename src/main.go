package main

import (
	"github.com/ITLab-Projects/service/api/v1"
	"github.com/ITLab-Projects/pkg/config"
	"github.com/ITLab-Projects/app"
)

func main() {
	cfg := config.GetConfig()
	app := app.New(cfg)
	app.AddApi(
		v1.New(
			app.Repository,
			app.Requester,
		),
	)

	app.Start()
}