package server

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ditatompel/xmr-remote-nodes/internal/database"
	"github.com/ditatompel/xmr-remote-nodes/internal/monero"

	"github.com/spf13/cobra"
)

var probersCmd = &cobra.Command{
	Use:   "probers",
	Short: "Add, edit, delete, and show registered probers",
	Long: `Command to administer prober machines.

This command should only be run on the server which directly connect to the MySQL database.
	`,
	Run: func(cmd *cobra.Command, _ []string) {
		if err := cmd.Help(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	},
}

var listProbersCmd = &cobra.Command{
	Use:   "list [search]",
	Short: "Print registered probers",
	Long: `Print list of registered prober machines.

Use [search] args to filter results by name or api key.

"sort-by" flag can be "id" or "last_submit_ts"
"sort-dir" flag can be "asc" or "desc"`,
	Example: `# To sort probers by last submit time in ascending order that contains "sin1":
xmr-nodes probers list -s last_submit_ts -d asc sin1`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := database.ConnectDB(); err != nil {
			panic(err)
		}
		sortBy, _ := cmd.Flags().GetString("sort-by")
		sortDir, _ := cmd.Flags().GetString("sort-dir")

		probersRepo := monero.NewProber()
		probers, err := probersRepo.Probers(monero.QueryProbers{
			Search:        strings.Join(args, " "),
			SortBy:        sortBy,
			SortDirection: sortDir,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(probers) == 0 {
			fmt.Println("No probers found")
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintf(w, "ID\t| Name\t| Last Submit\t| API Key\n")
		for _, prober := range probers {
			fmt.Fprintf(w, "%d\t| %s\t| %s\t| %s\n",
				prober.ID,
				prober.Name,
				time.Unix(prober.LastSubmitTS, 0).Format(time.RFC3339),
				prober.APIKey,
			)
		}
		w.Flush()
	},
}

func stringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

var addProbersCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add new prober",
	Long: `Create new prober identified by [name] (if provided) along with an API key.

This command will display the prober name and API key when successfully executed.`,
	Run: func(_ *cobra.Command, args []string) {
		if err := database.ConnectDB(); err != nil {
			fmt.Println(err)
			return
		}

		proberName := ""
		if len(args) > 0 {
			proberName = strings.Join(args, " ")
		} else {
			proberName = stringPrompt("Prober Name:")
		}

		proberRepo := monero.NewProber()
		prober, err := proberRepo.Add(proberName)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Name: %s\nAPI Key: %s\n", prober.Name, prober.APIKey)
	},
}

var editProbersCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit prober",
	Long:  `Edit prober name by id.`,
	Run: func(cmd *cobra.Command, _ []string) {
		if err := database.ConnectDB(); err != nil {
			fmt.Println(err)
			return
		}
		proberId, err := strconv.Atoi(stringPrompt("Prober ID:"))
		if err != nil {
			fmt.Println("Invalid ID:", err)
			return
		}

		proberName := stringPrompt("Prober Name:")
		proberRepo := monero.NewProber()
		err = proberRepo.Edit(proberId, proberName)
		if err != nil {
			fmt.Println("Failed to update prober:", err)
			return
		}

		fmt.Printf("Prober ID %d updated\n", proberId)
	},
}

var deleteProbersCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete prober",
	Long:  `Delete prober identified by id.`,
	Run: func(_ *cobra.Command, _ []string) {
		if err := database.ConnectDB(); err != nil {
			fmt.Println(err)
			return
		}
		proberId, err := strconv.Atoi(stringPrompt("Prober ID:"))
		if err != nil {
			fmt.Println("Invalid ID:", err)
			return
		}

		proberRepo := monero.NewProber()
		err = proberRepo.Delete(proberId)
		if err != nil {
			fmt.Println("Failed to delete prober:", err)
			return
		}

		fmt.Printf("Prober ID %d deleted\n", proberId)
	},
}
