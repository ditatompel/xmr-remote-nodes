package handler

import (
	"github.com/gofiber/fiber/v2"
)

func AppRoute(app *fiber.App) {
	app.Post("/auth/login", Login)
	app.Post("/auth/logout", Logout)
}

func V1Api(app *fiber.App) {
	v1 := app.Group("/api/v1")

	v1.Get("/prober", Prober)
	v1.Post("/prober", Prober)
	v1.Get("/crons", Crons)
}
