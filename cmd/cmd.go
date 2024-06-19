package cmd

import (
	"os"
	"xmr-remote-nodes/cmd/client"
	"xmr-remote-nodes/internal/config"

	"github.com/spf13/cobra"
)

const AppVer = "0.0.1"

var configFile string

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
	cobra.OnInitialize(initConfig)
	Root.PersistentFlags().StringVarP(&configFile, "config-file", "c", "", "Default to .env")
	Root.AddCommand(client.ProbeCmd)
	client.ProbeCmd.Flags().StringP("endpoint", "e", "", "Server endpoint")
	client.ProbeCmd.Flags().Bool("no-tor", false, "Only probe clearnet nodes")
}

func initConfig() {
	if configFile != "" {
		config.LoadAll(configFile)
		return
	}
	config.LoadAll(".env")
}
