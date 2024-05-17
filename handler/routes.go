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

	v1.Get("/crons", CookieProtected, Crons)
	v1.Get("/nodes", MoneroNodes)
	v1.Post("/nodes", AddNode)
	v1.Get("/nodes/id/:id", MoneroNode)
	v1.Get("/nodes/logs", ProbeLogs)
	v1.Get("/fees", NetFee)
	v1.Get("/countries", Countries)
	v1.Get("/job", CheckProber, GiveJob)
	v1.Post("/job", CheckProber, ProcessJob)
}
