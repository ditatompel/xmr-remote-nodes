package handler

func (s *fiberServer) Routes() {
	s.App.Get("/", s.homeHandler)
	s.App.Get("/remote-nodes", s.remoteNodesHandler)
	s.App.Get("/remote-nodes/id/:id", s.nodeHandler)
	s.App.Get("/add-node", s.addNodeHandler)
	s.App.Put("/add-node", s.addNodeHandler)

	// This is temporary route to redirect old path to new one. Once search
	// engine results updated to the new path, this route should be removed.
	s.App.Get("/remote-nodes/logs", s.redirectLogs)

	// V1 API routes
	v1 := s.App.Group("/api/v1")

	// these routes are public, they don't require a prober api key
	v1.Get("/nodes", Nodes)
	v1.Post("/nodes", AddNode) // old add node form action endpoint. Deprecated: Use PUT /add-node instead
	v1.Get("/nodes/id/:id", Node)
	v1.Get("/nodes/logs", ProbeLogs)
	v1.Get("/fees", NetFees)
	v1.Get("/countries", Countries)

	// these routes are for prober, they require a prober api key
	v1.Get("/job", CheckProber, GiveJob)
	v1.Post("/job", CheckProber, ProcessJob)
}
