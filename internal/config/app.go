package config

import (
	"log/slog"
	"os"
	"strconv"
)

var Version string

type App struct {
	// general config
	LogLevel string

	// configuration for server
	URL    string // URL where user can access the web UI, don't put trailing slash
	Secret string // random 64-character hex string that give us 32 random bytes

	// fiber specific config
	Prefork     bool
	Host        string
	Port        int
	ProxyHeader string
	AllowOrigin string

	// configuration for prober (client)
	ServerEndpoint string
	APIKey         string
	AcceptTor      bool
	TorSOCKS       string
	AcceptI2P      bool
	I2PSOCKS       string
	IPv6Capable    bool
}

func init() {
	if Version == "" {
		Version = "v0.0.0-alpha.0.000000.dev"
	}
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
	app.URL = os.Getenv("APP_URL")
	app.Secret = os.Getenv("APP_SECRET")

	// fiber specific config
	app.Host = os.Getenv("APP_HOST")
	app.Port, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	app.Prefork, _ = strconv.ParseBool(os.Getenv("APP_PREFORK"))
	app.ProxyHeader = os.Getenv("APP_PROXY_HEADER")
	app.AllowOrigin = os.Getenv("APP_ALLOW_ORIGIN")

	// prober configuration
	app.ServerEndpoint = os.Getenv("SERVER_ENDPOINT")
	app.APIKey = os.Getenv("API_KEY")
	app.AcceptTor, _ = strconv.ParseBool(os.Getenv("ACCEPT_TOR"))
	app.TorSOCKS = os.Getenv("TOR_SOCKS")
	app.AcceptI2P, _ = strconv.ParseBool(os.Getenv("ACCEPT_I2P"))
	app.I2PSOCKS = os.Getenv("I2P_SOCKS")
	app.IPv6Capable, _ = strconv.ParseBool(os.Getenv("IPV6_CAPABLE"))
}
