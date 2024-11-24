package handler

import (
	"github.com/ditatompel/xmr-remote-nodes/internal/config"
	"github.com/ditatompel/xmr-remote-nodes/internal/database"
	"github.com/gofiber/fiber/v2"
)

type fiberServer struct {
	*fiber.App
	db     *database.DB
	url    string
	secret string
}

// NewServer returns a new fiber server
func NewServer() *fiberServer {
	if database.ConnectDB() != nil {
		panic("Failed to connect to database")
	}
	server := &fiberServer{
		App: fiber.New(fiber.Config{
			Prefork:     config.AppCfg().Prefork,
			ProxyHeader: config.AppCfg().ProxyHeader,
			AppName:     "XMR Nodes Aggregator " + config.Version,
		}),
		db:  database.GetDB(),
		url: config.AppCfg().URL,
	}

	return server
}
