package handler

import (
	"github.com/gofiber/fiber/v2"
)

// V1 API routes
func V1Api(app *fiber.App) {
	v1 := app.Group("/api/v1")

	// these routes are public, they don't require a prober api key
	v1.Get("/nodes", Nodes)
	v1.Post("/nodes", AddNode)
	v1.Get("/nodes/id/:id", Node)
	v1.Get("/nodes/logs", ProbeLogs)
	v1.Get("/fees", NetFees)
	v1.Get("/countries", Countries)

	// these routes are for prober, they require a prober api key
	v1.Get("/job", CheckProber, GiveJob)
	v1.Post("/job", CheckProber, ProcessJob)
}
