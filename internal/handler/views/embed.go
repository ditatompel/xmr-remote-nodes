package views

import (
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed assets/*
var embedStatic embed.FS

func EmbedAssets() fiber.Handler {
	return filesystem.New(filesystem.Config{
		Root:       http.FS(embedStatic),
		PathPrefix: "assets",
		Browse:     false,
	})
}
