package handler

import (
	"fmt"
	"strconv"

	"github.com/a-h/templ"
	"github.com/ditatompel/xmr-remote-nodes/internal/handler/views"
	"github.com/ditatompel/xmr-remote-nodes/internal/monero"
	"github.com/ditatompel/xmr-remote-nodes/internal/paging"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

// Redirect old `/remote-nodes/logs/?node_id={id}` path to `/remote-nodes/id/{id}`
//
// This is temporary handler to redirect old path to new one. Once search
// engine results updated to the new path, this handler should be removed.
func (s *fiberServer) redirectLogs(c *fiber.Ctx) error {
	id := c.QueryInt("node_id", 0)
	if id == 0 {
		return c.Redirect("/remote-nodes", fiber.StatusMovedPermanently)
	}

	return c.Redirect(fmt.Sprintf("/remote-nodes/id/%d", id), fiber.StatusMovedPermanently)
}

// Render robots.txt
func (s *fiberServer) robotsTxtHandler(c *fiber.Ctx) error {
	return c.SendString("User-agent: *\nAllow: /\n")
}

// Render Home Page
func (s *fiberServer) homeHandler(c *fiber.Ctx) error {
	p := views.Meta{
		Title:       "Monero Remote Node",
		Description: "A website that helps you monitor your favourite Monero remote nodes, but YOU BETTER RUN AND USE YOUR OWN NODE.",
		Keywords:    "monero,monero,xmr,monero node,xmrnode,cryptocurrency,monero remote node,monero testnet,monero stagenet",
		Robots:      "INDEX,FOLLOW",
		Permalink:   s.url,
		Identifier:  "/",
	}

	c.Set("Link", fmt.Sprintf(`<%s>; rel="canonical"`, p.Permalink))
	home := views.BaseLayout(p, views.Home())
	handler := adaptor.HTTPHandler(templ.Handler(home))

	return handler(c)
}

// Render Add Node Page
func (s *fiberServer) addNodeHandler(c *fiber.Ctx) error {
	switch c.Method() {
	case fiber.MethodPut:
		type formData struct {
			Protocol string `form:"protocol"`
			Hostname string `form:"hostname"`
			Port     int    `form:"port"`
		}
		var f formData

		if err := c.BodyParser(&f); err != nil {
			handler := adaptor.HTTPHandler(templ.Handler(views.Alert("error", "Cannot parse the request body")))
			return handler(c)
		}

		moneroRepo := monero.New()
		if err := moneroRepo.Add(c.IP(), s.secret, f.Protocol, f.Hostname, uint(f.Port)); err != nil {
			handler := adaptor.HTTPHandler(templ.Handler(views.Alert("error", err.Error())))
			return handler(c)
		}

		handler := adaptor.HTTPHandler(templ.Handler(views.Alert("success", "Node added successfully")))
		return handler(c)
	}
	p := views.Meta{
		Title:       "Add Monero Node",
		Description: "You can use this page to add known remote node to the system so my bots can monitor it.",
		Keywords:    "monero,monero node,monero public node,monero wallet,list monero node,monero node monitoring",
		Robots:      "INDEX,FOLLOW",
		Permalink:   s.url + "/add-node",
		Identifier:  "/add-node",
	}

	c.Set("Link", fmt.Sprintf(`<%s>; rel="canonical"`, p.Permalink))
	home := views.BaseLayout(p, views.AddNode())
	handler := adaptor.HTTPHandler(templ.Handler(home))

	return handler(c)
}

// Returns a single node information based on `id` query param (API endpoint, JSON data)
func (s *fiberServer) nodeAPI(c *fiber.Ctx) error {
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

// Render Remote Nodes Page
func (s *fiberServer) remoteNodesHandler(c *fiber.Ctx) error {
	p := views.Meta{
		Title:       "Public Monero Remote Nodes List",
		Description: "Although it's possible to use these existing public Monero nodes, you're MUST RUN AND USE YOUR OWN NODE!",
		Keywords:    "monero remote nodes,public monero nodes,monero public nodes,monero wallet,tor monero node,monero cors rpc",
		Robots:      "INDEX,FOLLOW",
		Permalink:   s.url + "/remote-nodes",
		Identifier:  "/remote-nodes",
	}

	moneroRepo := monero.New()
	query := monero.QueryNodes{
		Paging: paging.Paging{
			Limit:         c.QueryInt("limit", 10), // rows per page
			Page:          c.QueryInt("page", 1),
			SortBy:        c.Query("sort_by", "last_checked"),
			SortDirection: c.Query("sort_direction", "desc"),
			Refresh:       c.Query("refresh"),
		},
		Host:     c.Query("host"),
		Nettype:  c.Query("nettype", "any"),
		Protocol: c.Query("protocol", "any"),
		CC:       c.Query("cc", "any"),
		Status:   c.QueryInt("status", -1),
		CORS:     c.Query("cors"),
	}

	nodes, err := moneroRepo.Nodes(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	countries, err := moneroRepo.Countries()
	if err != nil {
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	pagination := paging.NewPagination(query.Page, nodes.TotalPages)

	// handle request from HTMX
	if c.Get("HX-Target") == "tbl_nodes" {
		cmp := views.BlankLayout(views.TableNodes(nodes, countries, query, pagination))
		handler := adaptor.HTTPHandler(templ.Handler(cmp))
		return handler(c)
	}

	c.Set("Link", fmt.Sprintf(`<%s>; rel="canonical"`, p.Permalink))
	home := views.BaseLayout(p, views.RemoteNodes(nodes, countries, query, pagination))
	handler := adaptor.HTTPHandler(templ.Handler(home))

	return handler(c)
}

// Returns a single node information based on `id` query param.
// This used for node modal and node details page including node probe logs.
func (s *fiberServer) nodeHandler(c *fiber.Ctx) error {
	nodeID, err := c.ParamsInt("id", 0)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	if nodeID == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid node id",
			"data":    nil,
		})
	}

	moneroRepo := monero.New()
	node, err := moneroRepo.Node(nodeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	switch c.Get("HX-Target") {
	case "modal-section":
		cmp := views.ModalLayout(fmt.Sprintf("Node #%d", nodeID), views.Node(node))
		handler := adaptor.HTTPHandler(templ.Handler(cmp))
		return handler(c)
	}

	queryLogs := monero.QueryLogs{
		Paging: paging.Paging{
			Limit:         c.QueryInt("limit", 10), // rows per page
			Page:          c.QueryInt("page", 1),
			SortBy:        c.Query("sort_by", "id"),
			SortDirection: c.Query("sort_direction", "desc"),
			Refresh:       c.Query("refresh"),
		},
		NodeID:       int(node.ID),
		Status:       c.QueryInt("status", -1),
		FailedReason: c.Query("failed_reason"),
	}

	logs, err := moneroRepo.Logs(queryLogs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	pagination := paging.NewPagination(queryLogs.Page, logs.TotalPages)

	// handle datatable logs filters, sort request from HTMX
	if c.Get("HX-Target") == "tbl_logs" {
		cmp := views.BlankLayout(views.TableLogs(fmt.Sprintf("/remote-nodes/id/%d", node.ID), logs, queryLogs, pagination))
		handler := adaptor.HTTPHandler(templ.Handler(cmp))
		return handler(c)
	}

	p := views.Meta{
		Title:       fmt.Sprintf("%s on Port %d", node.Hostname, node.Port),
		Description: fmt.Sprintf("Monero %s remote node %s running on port %d", node.Nettype, node.Hostname, node.Port),
		Keywords:    fmt.Sprintf("monero log,monero node log,monitoring monero log,monero,xmr,monero node,xmrnode,cryptocurrency,monero %s,%s", node.Nettype, node.Hostname),
		Robots:      "INDEX,FOLLOW",
		Permalink:   s.url + "/remote-nodes/id/" + strconv.Itoa(int(node.ID)),
		Identifier:  "/remote-nodes",
	}

	c.Set("Link", fmt.Sprintf(`<%s>; rel="canonical"`, p.Permalink))
	cmp := views.BaseLayout(p, views.NodeDetails(node, logs, queryLogs, pagination))
	handler := adaptor.HTTPHandler(templ.Handler(cmp))
	return handler(c)
}

// Returns list of nodes (API endpoint, JSON data)
func (s *fiberServer) nodesAPI(c *fiber.Ctx) error {
	moneroRepo := monero.New()
	query := monero.QueryNodes{
		Paging: paging.Paging{
			Limit:         c.QueryInt("limit", 10), // rows per page
			Page:          c.QueryInt("page", 1),
			SortBy:        c.Query("sort_by", "last_checked"),
			SortDirection: c.Query("sort_direction", "desc"),
			Refresh:       c.Query("refresh"),
		},
		Host:     c.Query("host"),
		Nettype:  c.Query("nettype", "any"),
		Protocol: c.Query("protocol", "any"),
		CC:       c.Query("cc", "any"),
		Status:   c.QueryInt("status", -1),
		CORS:     c.Query("cors"),
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

// Returns probe logs reported by nodes (API endpoint, JSON data)
func (s *fiberServer) probeLogsAPI(c *fiber.Ctx) error {
	moneroRepo := monero.New()
	query := monero.QueryLogs{
		Paging: paging.Paging{
			Limit:         c.QueryInt("limit", 10), // rows per page
			Page:          c.QueryInt("page", 1),
			SortBy:        c.Query("sort_by", "id"),
			SortDirection: c.Query("sort_direction", "desc"),
			Refresh:       c.Query("refresh"),
		},
		NodeID:       c.QueryInt("node_id", 0),
		Status:       c.QueryInt("status", -1),
		FailedReason: c.Query("failed_reason"),
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

// Handles `POST /nodes` request to add a new node
//
// Deprecated: AddNode is deprecated, use s.addNodeHandler with put method instead
func (s *fiberServer) addNodeAPI(c *fiber.Ctx) error {
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
	if err := moneroRepo.Add(c.IP(), s.secret, protocol, hostname, uint(port)); err != nil {
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

// Returns majority network fees (API endpoint, JSON data)
func (s *fiberServer) netFeesAPI(c *fiber.Ctx) error {
	moneroRepo := monero.New()
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Success",
		"data":    moneroRepo.NetFees(),
	})
}

// Returns list of countries, count by nodes (API endpoint, JSON data)
func (s *fiberServer) countriesAPI(c *fiber.Ctx) error {
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

// Returns node to be probed by the prober (API endpoint, JSON data)
//
// This handler should protected by `s.checkProberMW` middleware.
func (s *fiberServer) giveJobAPI(c *fiber.Ctx) error {
	acceptTor := c.QueryInt("accept_tor", 0)
	acceptI2P := c.QueryInt("accept_i2p", 0)
	acceptIPv6 := c.QueryInt("accept_ipv6", 0)

	moneroRepo := monero.New()
	node, err := moneroRepo.GiveJob(acceptTor, acceptI2P, acceptIPv6)
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

// Handles probe report submission by the prober (API endpoint, JSON data)
//
// This handler should protected by `CheckProber` middleware.
func (s *fiberServer) processJobAPI(c *fiber.Ctx) error {
	var report monero.ProbeReport

	if err := c.BodyParser(&report); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	moneroRepo := monero.New()

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
