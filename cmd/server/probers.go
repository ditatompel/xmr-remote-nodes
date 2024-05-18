package server

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
	"xmr-remote-nodes/cmd"
	"xmr-remote-nodes/internal/database"
	"xmr-remote-nodes/internal/repo"

	"github.com/spf13/cobra"
)

var probersCmd = &cobra.Command{
	Use:   "probers",
	Short: "[server] Add, edit, delete, and show registered probers",
	Long: `Command to administer prober machines.

This command should only be run on the server which directly connect to the MySQL database.
	`,
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
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

		probersRepo := repo.NewProberRepo(database.GetDB())
		probers, err := probersRepo.Probers(repo.ProbersQueryParams{
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
				prober.Id,
				prober.Name,
				time.Unix(prober.LastSubmitTs, 0).Format(time.RFC3339),
				prober.ApiKey,
			)
		}
		w.Flush()
	},
}

var addProbersCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "[server] Add new prober",
	Long:  `Add new prober machine.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: Add new prober")
	},
}

func init() {
	cmd.Root.AddCommand(probersCmd)
	probersCmd.AddCommand(listProbersCmd)
	probersCmd.AddCommand(addProbersCmd)
	listProbersCmd.Flags().StringP("sort-by", "s", "last_submit_ts", "Sort by column name, can be id or last_submit_ts")
	listProbersCmd.Flags().StringP("sort-dir", "d", "desc", "Sort direction, can be asc or desc")
}
