package main

import (
	"github.com/ITLab-Projects/service/api/v2"
	"github.com/ITLab-Projects/service/api/v1"
	"github.com/ITLab-Projects/pkg/config"
	"github.com/ITLab-Projects/app"
	
)

// @title ITLab-Projects API
// @version 1.0
// @description This is a server to get projects from github
// @BasePath /api/projects
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization


func main() {
	cfg := config.GetConfig()
	app := app.New(cfg)
	app.AddApi(
		v1.New(
			v1.Config{
				Testmode: cfg.App.TestMode,
				Config: *cfg.Auth,
				UpdateTime: cfg.App.UpdateTime,
			},
			app.Repository,
			app.Requester,
			app.MFSRequester,
		),
		v2.New(
			v2.Config{
				Testmode: cfg.App.TestMode,
				Config: *cfg.Auth,
			},
			app.Repository,
		),
	)

	app.Start()
}