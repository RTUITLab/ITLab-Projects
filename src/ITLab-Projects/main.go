package main

import (
	"github.com/ITLab-Projects/service/api/v1"
	"github.com/ITLab-Projects/pkg/config"
	"github.com/ITLab-Projects/app"
	
)

// @title ITLab-Projects API
// @version 1.0
// @description This is a server to get projects from github


func main() {
	cfg := config.GetConfig()
	app := app.New(cfg)
	app.AddApi(
		v1.New(
			v1.Config{
				cfg.App.TestMode,
				*cfg.Auth,
			},
			app.Repository,
			app.Requester,
			app.MFSRequester,
		).WithUpdater(cfg.App.UpdateTime),
	)

	app.Start()
}