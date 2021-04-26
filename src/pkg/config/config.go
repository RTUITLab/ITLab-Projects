package config

// Package config need to configure DataBase connection
// and to connect to github and some start up settings for main App

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB 			*DBConfig		`json:"DbOptions"`
	Auth 		*AuthConfig		`json:"AuthOptions"`
	App 		*AppConfig		`json:"AppOptions"`
	Services	*OtherServicesConfig
}

type OtherServicesConfig struct {
	MFS		*MFSConfig
}

type MFSConfig struct {
	BaseURL		string `envconfig:"ITLAB_PROJECTS_MFSURL"`
}

type DBConfig struct {
	URI 		string		`envconfig:"ITLAB_PROJECTS_DBURI" json:"uri"`
}
type AuthConfig struct {
	KeyURL		string		`envconfig:"ITLABPROJ_KEYURL" json:"keyUrl"`
	Audience	string		`envconfig:"ITLABPROJ_AUDIENCE" json:"audience"`
	Issuer		string		`envconfig:"ITLABPROJ_ISSUER" json:"issuer"`
	Scope		string		`envconfig:"ITLABPROJ_SCOPE" json:"scope"`
	Github		Github		`json:"Github"`
}

type Github struct {
	AppID			int64		`json:"appID"`
	PathToPem		string		`json:"pathToPem"`
	Installation 	string		`json:"installation"`
	AccessToken   	string		`envconfig:"ITLAB_PROJECTS_ACCESSKEY" json:"accessToken"`
}

type AppConfig struct {
	AppPort				string	`envconfig:"ITLAB_PROJECTS_APPPORT" json:"appPort"`
	TestMode			bool	`envconfig:"ITLAB_PROJECTS_TESTMODE" json:"testMode"`
}

func GetConfig() *Config {
	var config Config
	if err := godotenv.Load("./.env"); err != nil {
		panic(err)
	}

	if err := envconfig.Process("itlab_projects", &config); err != nil {
		log.WithFields(
			log.Fields{
				"function" : "envconfig.Process",
				"error"	:	err,
			},
		).Fatal("Can't read env vars, shutting down...")
	}
	return &config
}
