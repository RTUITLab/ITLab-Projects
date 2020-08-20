package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type Config struct {
	DB *DBConfig		`json:"DbOptions"`
	Auth *AuthConfig	`json:"AuthOptions"`
	App *AppConfig		`json:"AppOptions"`
}

type DBConfig struct {
	Host 					string		`json:"host"`
	DBPort 					string		`json:"dbPort"`
	DBName 					string		`json:"dbName"`
	ProjectsCollectionName	string		`json:"projectsCollectionName"`
	ReposCollectionName		string		`json:"reposCollectionName"`
}
type AuthConfig struct {
	KeyURL		string		`json:"keyUrl"`
	Audience	string		`json:"audience"`
	Issuer		string		`json:"issuer"`
	Scope		string		`json:"scope"`
	Github		Github		`json:"Github"`
	Gitlab		Gitlab		`json:"Gitlab"`
}
type Github struct {
	AppID		int64		`json:"appID"`
	PathToPem	string		`json:"pathToPem"`
	Installation string		`json:"installation"`
	AccessToken   string	`json:"accessToken"`
}
type Installation struct {
	ID			int64		`json:"id"`
	Account     Account		`json:"account"`
}
type Account struct {
	Login 		string		`json:"login"`
}
type Gitlab struct {
	AccessToken   string	`json:"accessToken"`
}
type AppConfig struct {
	AppPort		string	`json:"appPort"`
	TestMode	bool	`json:"testMode"`
	ElemsPerPage int	`json:"elemsPerPage"`
}

func GetConfig() *Config {
	var config Config
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "GetConfig.ReadFile",
			"error"	:	err,
		},
		).Fatal("Can't read config.json file, shutting down...")
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "GetConfig.Unmarshal",
			"error"	:	err,
		},
		).Fatal("Can't correctly parse json from config.json, shutting down...")
	}

	data, err = ioutil.ReadFile("auth_config.json")
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "GetConfig.ReadFile",
			"error"	:	err,
		},
		).Fatal("Can't read auth_config.json file, shutting down...")
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "GetConfig.Unmarshal",
			"error"	:	err,
		},
		).Fatal("Can't correctly parse json from auth_config.json, shutting down...")
	}
	return &config
}
