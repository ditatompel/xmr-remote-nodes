package handler

import (
	"xmr-remote-nodes/internal/database"
	"xmr-remote-nodes/internal/repo"

	"github.com/gofiber/fiber/v2"
)

func CheckProber(c *fiber.Ctx) error {
	key := c.Get("X-Prober-Api-Key")
	if key == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
			"data":    nil,
		})
	}

	proberRepo := repo.NewProberRepo(database.GetDB())

	prober, err := proberRepo.CheckApi(key)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "No API key match",
			"data":    nil,
		})
	}

	c.Locals("prober_id", prober.Id)
	return c.Next()
}
