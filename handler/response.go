package handler

import (
	"strconv"
	"xmr-remote-nodes/internal/monero"

	"github.com/gofiber/fiber/v2"
)

func MoneroNode(c *fiber.Ctx) error {
	nodeId, err := c.ParamsInt("id", 0)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	if nodeId == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid node id",
			"data":    nil,
		})
	}

	moneroRepo := monero.New()
	node, err := moneroRepo.Node(nodeId)
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
		"data":    node,
	})
}

func MoneroNodes(c *fiber.Ctx) error {
	moneroRepo := monero.New()
	query := monero.QueryNodes{
		RowsPerPage:   c.QueryInt("limit", 10),
		Page:          c.QueryInt("page", 1),
		SortBy:        c.Query("sort_by", "id"),
		SortDirection: c.Query("sort_direction", "desc"),
		Host:          c.Query("host"),
		Nettype:       c.Query("nettype", "any"),
		Protocol:      c.Query("protocol", "any"),
		CC:            c.Query("cc", "any"),
		Status:        c.QueryInt("status", -1),
		CORS:          c.QueryInt("cors", -1),
	}

	nodes, err := moneroRepo.Nodes(query)
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
		"data":    nodes,
	})
}

func ProbeLogs(c *fiber.Ctx) error {
	moneroRepo := monero.New()
	query := monero.QueryLogs{
		RowsPerPage:   c.QueryInt("limit", 10),
		Page:          c.QueryInt("page", 1),
		SortBy:        c.Query("sort_by", "id"),
		SortDirection: c.Query("sort_direction", "desc"),
		NodeID:        c.QueryInt("node_id", 0),
		Status:        c.QueryInt("status", -1),
		FailedReason:  c.Query("failed_reason"),
	}

	logs, err := moneroRepo.Logs(query)
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
		"data":    logs,
	})
}

func AddNode(c *fiber.Ctx) error {
	formPort := c.FormValue("port")
	port, err := strconv.Atoi(formPort)
	if err != nil {
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid port number",
			"data":    nil,
		})
	}

	protocol := c.FormValue("protocol")
	hostname := c.FormValue("hostname")

	moneroRepo := monero.New()
	if err := moneroRepo.Add(protocol, hostname, uint(port)); err != nil {
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Query Ok",
		"data":    nil,
	})
}

func NetFee(c *fiber.Ctx) error {
	moneroRepo := monero.New()
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Success",
		"data":    moneroRepo.NetFee(),
	})
}

func Countries(c *fiber.Ctx) error {
	moneroRepo := monero.New()
	countries, err := moneroRepo.Countries()
	if err != nil {
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Success",
		"data":    countries,
	})
}

func GiveJob(c *fiber.Ctx) error {
	acceptTor := c.QueryInt("accept_tor", 0)

	moneroRepo := monero.New()
	node, err := moneroRepo.GiveJob(acceptTor)
	if err != nil {
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Success",
		"data":    node,
	})
}

func ProcessJob(c *fiber.Ctx) error {
	moneroRepo := monero.New()
	report := monero.ProbeReport{}

	if err := c.BodyParser(&report); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	if err := moneroRepo.ProcessJob(report, c.Locals("prober_id").(int64)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Success",
		"data":    nil,
	})
}
