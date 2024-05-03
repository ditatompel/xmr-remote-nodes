package handler

import (
	"fmt"
	"time"

	"github.com/ditatompel/xmr-nodes/internal/database"
	"github.com/ditatompel/xmr-nodes/internal/repo"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	payload := repo.Admin{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	repo := repo.NewAdminRepo(database.GetDB())
	res, err := repo.Login(payload.Username, payload.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	token := fmt.Sprintf("auth_%d_%d", res.Id, time.Now().Unix())
	c.Cookie(&fiber.Cookie{
		Name:     "xmr-nodes-ui",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Logged in",
		"data":    nil,
	})
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "xmr-nodes-ui",
		Value:    "",
		Expires:  time.Now(),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Logged out",
		"data":    nil,
	})
}

func Prober(c *fiber.Ctx) error {
	proberRepo := repo.NewProberRepo(database.GetDB())

	if c.Method() == "POST" {
		payload := repo.Prober{}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
				"data":    nil,
			})
		}
		if payload.Name == "" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  "error",
				"message": "Please fill prober name",
				"data":    nil,
			})
		}
		err := proberRepo.AddProber(payload.Name)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
				"data":    nil,
			})
		}
	}

	query := repo.ProbersQueryParams{
		RowsPerPage:   c.QueryInt("limit", 10),
		Page:          c.QueryInt("page", 1),
		SortBy:        c.Query("sort_by", "id"),
		SortDirection: c.Query("sort_direction", "desc"),
		Name:          c.Query("name"),
		ApiKey:        c.Query("api_key"),
	}

	prober, err := proberRepo.Probers(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Success",
		"data":    prober,
	})
}

func Crons(c *fiber.Ctx) error {
	cronRepo := repo.NewCron(database.GetDB())

	crons, err := cronRepo.Crons()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Crons",
		"data":    crons,
	})
}
