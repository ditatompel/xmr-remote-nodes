package server

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
	"xmr-remote-nodes/internal/database"
	"xmr-remote-nodes/internal/repo"

	"github.com/spf13/cobra"
)

var probersCmd = &cobra.Command{
	Use:   "probers",
	Short: "Add, edit, delete, and show registered probers",
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
	Short: "Add new prober",
	Long: `Create new prober identified by [name] (if provided) along with an API key.

This command will display the prober name and API key when successfully executed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := database.ConnectDB(); err != nil {
			panic(err)
		}

		proberName := ""
		if len(args) > 0 {
			proberName = strings.Join(args, " ")
		} else {
			proberName = stringPrompt("Prober Name:")
		}

		proberRepo := repo.NewProberRepo(database.GetDB())
		prober, err := proberRepo.Add(proberName)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Name: %s\nAPI Key: %s\n", prober.Name, prober.ApiKey)
	},
}
