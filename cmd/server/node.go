package server

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/ditatompel/xmr-remote-nodes/internal/database"
	"github.com/ditatompel/xmr-remote-nodes/internal/monero"

	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "[Server] Administer monitored nodes",
	Long: `Command to administer monitored nodes.

This command should only be run on the server which directly connect to the MySQL database.
	`,
	Run: func(cmd *cobra.Command, _ []string) {
		if err := cmd.Help(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}

var deleteNodeCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete node",
	Long: `Delete node identified by ID.

This command  delete node and it's associated probe logs (if exists).

To find out the node ID, visit frontend UI or from "/api/v1/nodes" endpoint.
	`,
	Run: func(_ *cobra.Command, _ []string) {
		if err := database.ConnectDB(); err != nil {
			fmt.Println(err)
			return
		}
		nodeID, err := strconv.Atoi(stringPrompt("Node ID:"))
		if err != nil {
			fmt.Println("Invalid ID:", err)
			return
		}

		moneroRepo := monero.New()
		err = moneroRepo.Delete(uint(nodeID))
		if err != nil {
			fmt.Println("Failed to delete node:", err)
			return
		}

		fmt.Printf("Node ID %d deleted\n", nodeID)
	},
}
