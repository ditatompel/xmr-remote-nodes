package handler

import (
	"github.com/ditatompel/xmr-remote-nodes/internal/monero"

	"github.com/gofiber/fiber/v2"
)

func CheckProber(c *fiber.Ctx) error {
	key := c.Get(monero.ProberAPIKey)
	if key == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
			"data":    nil,
		})
	}

	prober, err := monero.NewProber().CheckAPI(key)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "No API key match",
			"data":    nil,
		})
	}

	c.Locals("prober_id", prober.ID)
	return c.Next()
}
