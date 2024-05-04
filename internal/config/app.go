package config

import (
	"os"
	"strconv"
)

type App struct {
	// configuration for server
	Debug       bool
	Prefork     bool
	Host        string
	Port        int
	ProxyHeader string
	AllowOrigin string
	SecretKey   string
	LogLevel    string
	// configuration for prober (client)
	ServerEndpoint string
	ApiKey         string
	AcceptTor      bool
	TorSocks       string
}

var app = &App{}

func AppCfg() *App {
	return app
}

// loads App configuration
func LoadApp() {
	// server configuration
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
	// prober configuration
	app.ServerEndpoint = os.Getenv("SERVER_ENDPOINT")
	app.ApiKey = os.Getenv("API_KEY")
	app.AcceptTor, _ = strconv.ParseBool(os.Getenv("ACCEPT_TOR"))
	app.TorSocks = os.Getenv("TOR_SOCKS")
}
