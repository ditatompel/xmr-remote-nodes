package cmd

import (
	"os"
	"xmr-remote-nodes/cmd/client"

	"github.com/spf13/cobra"
)

const AppVer = "0.0.1"

var Root = &cobra.Command{
	Use:     "xmr-nodes",
	Short:   "XMR Nodes",
	Version: AppVer,
}

func Execute() {
	err := Root.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	Root.AddCommand(client.ProbeCmd)
}
