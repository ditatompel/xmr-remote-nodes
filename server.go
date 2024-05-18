//go:build server

package main

import (
	"xmr-remote-nodes/cmd"
	"xmr-remote-nodes/cmd/server"
)

func init() {
	cmd.Root.AddCommand(server.AdminCmd)
}
