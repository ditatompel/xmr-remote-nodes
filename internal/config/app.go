package config

import (
	"os"
	"strconv"
)

type App struct {
	Debug       bool
	Prefork     bool
	Host        string
	Port        int
	ProxyHeader string
	AllowOrigin string
	SecretKey   string
	LogLevel    string
}

var app = &App{}

func AppCfg() *App {
	return app
}

// LoadApp loads App configuration
func LoadApp() {
	app.Host = os.Getenv("APP_HOST")
	app.Port, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	app.Debug, _ = strconv.ParseBool(os.Getenv("APP_DEBUG"))
	app.Prefork, _ = strconv.ParseBool(os.Getenv("APP_PREFORK"))
	app.ProxyHeader = os.Getenv("APP_PROXY_HEADER")
	app.AllowOrigin = os.Getenv("APP_ALLOW_ORIGIN")
	app.SecretKey = os.Getenv("SECRET_KEY")
	app.LogLevel = os.Getenv("LOG_LEVEL")
	if app.LogLevel == "" {
		app.LogLevel = "INFO"
	}
	if app.Debug {
		app.LogLevel = "DEBUG"
	}
}
