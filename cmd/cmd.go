package cmd

import (
	"os"

	"github.com/ditatompel/xmr-remote-nodes/cmd/client"
	"github.com/ditatompel/xmr-remote-nodes/internal/config"

	"github.com/spf13/cobra"
)

var configFile string

var Root = &cobra.Command{
	Use:     "xmr-nodes",
	Short:   "XMR Nodes",
	Version: config.Version,
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
	client.ProbeCmd.Flags().Bool("no-tor", false, "Do not probe tor nodes")
	client.ProbeCmd.Flags().Bool("no-i2p", false, "Do not probe i2p nodes")
}

func initConfig() {
	if configFile != "" {
		config.LoadAll(configFile)
		return
	}
	config.LoadAll(".env")
}
