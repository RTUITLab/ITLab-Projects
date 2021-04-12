package config

// Package config need to configure DataBase connection
// and to connect to github and some start up settings for main App

import (
	"encoding/json"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB 		*DBConfig		`json:"DbOptions"`
	Auth 	*AuthConfig		`json:"AuthOptions"`
	App 	*AppConfig		`json:"AppOptions"`
}

type DBConfig struct {
	URI 		string		`envconfig:"ITLABPROJ_URI",json:"uri"`
}
type AuthConfig struct {
	KeyURL		string		`envconfig:"ITLABPROJ_KEYURL",json:"keyUrl"`
	Audience	string		`envconfig:"ITLABPROJ_AUDIENCE",json:"audience"`
	Issuer		string		`envconfig:"ITLABPROJ_ISSUER",json:"issuer"`
	Scope		string		`envconfig:"ITLABPROJ_SCOPE",json:"scope"`
	Github		Github		`json:"Github"`
	Gitlab		Gitlab		`json:"Gitlab"`
}
type Github struct {
	AppID			int64		`json:"appID"`
	PathToPem		string		`json:"pathToPem"`
	Installation 	string		`json:"installation"`
	AccessToken   	string		`envconfig:"ITLABPROJ_GHACCESSTOKEN",json:"accessToken"`
}
type Installation struct {
	ID				int64		`json:"id"`
	Account     	Account		`json:"account"`
}
type Account struct {
	Login 			string		`json:"login"`
}
type Gitlab struct {
	AccessToken   	string	`json:"accessToken"`
}
type AppConfig struct {
	AppPort				string	`envconfig:"ITLABPROJ_APPPORT",json:"appPort"`
	TestMode			bool	`envconfig:"ITLABPROJ_TESTMODE",json:"testMode"`
	ElemsPerPage 		int		`envconfig:"ITLABPROJ_ELEMSPERPAGE",json:"elemsPerPage"`
	ProjectFileBranch 	string	`envconfig:"ITLABPROJ_PROJFILEBRANCH",json:"projectFileBranch"`
}

func GetConfig() *Config {
	var config Config
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.WithFields(
			log.Fields{
				"function" : "GetConfig.ReadFile",
				"error"	:	err,
			},
		).Warning("Can't read config.json file, shutting down...")
	}
	if err = json.Unmarshal(data, &config); err != nil {
		log.WithFields(
			log.Fields{
				"function" : "GetConfig.Unmarshal",
				"error"	:	err,
			},
		).Warning("Can't correctly parse json from config.json, shutting down...")
	}

	data, err = ioutil.ReadFile("auth_config.json")
	if err != nil {
		log.WithFields(
			log.Fields{
				"function" : "GetConfig.ReadFile",
				"error"	:	err,
			},
		).Warning("Can't read auth_config.json file, shutting down...")
	}
	if err = json.Unmarshal(data, &config); err != nil {
		log.WithFields(
			log.Fields{
				"function" : "GetConfig.Unmarshal",
				"error"	:	err,
			},
		).Warning("Can't correctly parse json from auth_config.json, shutting down...")
	}

	if err = envconfig.Process("itlabproj", &config); err != nil {
		log.WithFields(
			log.Fields{
				"function" : "envconfig.Process",
				"error"	:	err,
			},
		).Fatal("Can't read env vars, shutting down...")
	}
	return &config
}
