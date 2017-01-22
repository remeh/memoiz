// Memoiz Backend
//
// Runtime configuration.
//
// Rémy Mathieu © 2016

package config

import (
	"strings"

	"remy.io/memoiz/log"

	"github.com/vrischmann/envconfig"
)

var Config Configuration

func init() {
	if c, err := Read(); err != nil {
		log.Error("while reading Config:", err.Error())
	} else {
		Config = c
	}
}

// ----------------------

type Configuration struct {
	// Address to listen to.
	ListenAddr string `envconfig:"ADDR,default=:8080"`
	// Public directory with pages and assets.
	PublicDir string `envconfig:"PUBLIC,default=public/dist"`
	// Resources directory
	ResDir string `envconfig:"RES,default=resources/"`
	// Connection string
	ConnString string `envconfig:"CONN,default=host=/var/run/postgresql sslmode=disable user=memoiz dbname=memoiz password=memoiz"`
	// Azure Web search api key
	BingApiKey string `envconfig:"BING,optional"`
	// Google Knowledge Graph api key
	KgApiKey string `envconfig:"KG,optional"`

	// URL of the app calling the backend. For CORS purpose.
	AppUrl string `envconfig:"APP_URL,default=http://app.memoiz.com"`

	// Smtp Configuration
	SmtpHost     string `envconfig:"SMTP_HOST,optional"`
	SmtpPort     int    `envconfig:"SMTP_PORT,optional"`
	SmtpLogin    string `envconfig:"SMTP_LOGIN,optional"`
	SmtpPassword string `envconfig:"SMTP_PWD,optional"`
}

// readConfig reads in the environment var
// the configuration to start the runtime.
func Read() (Configuration, error) {
	var c Configuration
	err := envconfig.Init(&c)

	if err != nil {
		return c, err
	}

	if !strings.HasSuffix(c.PublicDir, "/") {
		c.PublicDir += "/"
	}

	if !strings.HasSuffix(c.ResDir, "/") {
		c.ResDir += "/"
	}

	return c, nil
}
