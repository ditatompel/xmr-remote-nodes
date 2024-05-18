package main

import (
	"xmr-remote-nodes/cmd"
	"xmr-remote-nodes/internal/config"
)

func main() {
	config.LoadAll(".env")
	cmd.Execute()
}
