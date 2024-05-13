package cmd

import (
	"os"
	"xmr-remote-nodes/internal/config"

	"github.com/spf13/cobra"
)

const AppVer = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "xmr-nodes",
	Short:   "XMR Nodes",
	Version: AppVer,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	config.LoadAll(".env")
}
