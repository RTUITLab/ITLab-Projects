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
	KeyURL		string		`envconfig:"ITLAB_PROJECTS_KEYURL" json:"keyUrl"`
	Audience	string		`envconfig:"ITLAB_PROJECTS_AUDIENCE" json:"audience"`
	Issuer		string		`envconfig:"ITLAB_PROJECTS_ISSUER" json:"issuer"`
	Scope		string		`envconfig:"ITLAB_PROJECTS_SCOPE" json:"scope"`
	*RolesConfig
	Github		Github		`json:"Github"`
}

type RolesConfig struct {
	// looks like roles = "admin user" parse to ["admin", "user"]
	Roles string			`envconfig:"ITLAB_PROJECTS_ROLES" json:"roles"`
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
	UpdateTime			string	`envconfig:"ITLAB_PROJECTS_UPDATETIME" json:"update_time"`
}

func GetConfig() *Config {
	var config Config
	if err := godotenv.Load("./.env"); err != nil {
		log.Warn("Don't find .env file")
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
