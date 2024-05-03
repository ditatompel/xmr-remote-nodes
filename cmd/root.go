package cmd

import (
	"os"

	"github.com/ditatompel/xmr-nodes/internal/config"

	"github.com/spf13/cobra"
)

const AppVer = "0.0.1"

var LogLevel string

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
	LogLevel = config.AppCfg().LogLevel
}
