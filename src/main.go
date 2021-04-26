package main

import (
	_ "github.com/ITLab-Projects/docs"
	httpSwager "github.com/swaggo/http-swagger"
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
			app.Repository,
			app.Requester,
			app.MFSRequester,
		),
	)

	app.Router.PathPrefix("/swagger").Handler(httpSwager.WrapHandler)

	app.Start()
}