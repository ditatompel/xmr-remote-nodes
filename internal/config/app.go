package config

import (
	"log/slog"
	"os"
	"strconv"
)

type App struct {
	// general config
	LogLevel string

	// configuration for server
	Prefork     bool
	Host        string
	Port        int
	ProxyHeader string
	AllowOrigin string

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
	// general config
	app.LogLevel = os.Getenv("LOG_LEVEL")
	switch app.LogLevel {
	case "DEBUG":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "ERROR":
		slog.SetLogLoggerLevel(slog.LevelError)
	case "WARN":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	// server configuration
	app.Host = os.Getenv("APP_HOST")
	app.Port, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	app.Prefork, _ = strconv.ParseBool(os.Getenv("APP_PREFORK"))
	app.ProxyHeader = os.Getenv("APP_PROXY_HEADER")
	app.AllowOrigin = os.Getenv("APP_ALLOW_ORIGIN")

	// prober configuration
	app.ServerEndpoint = os.Getenv("SERVER_ENDPOINT")
	app.ApiKey = os.Getenv("API_KEY")
	app.AcceptTor, _ = strconv.ParseBool(os.Getenv("ACCEPT_TOR"))
	app.TorSocks = os.Getenv("TOR_SOCKS")
}
