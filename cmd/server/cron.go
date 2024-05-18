package server

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
	"xmr-remote-nodes/internal/database"
	"xmr-remote-nodes/internal/repo"

	"github.com/spf13/cobra"
)

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Print cron tasks",
	Long:  `Print list of regular cron tasks running on the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := database.ConnectDB(); err != nil {
			panic(err)
		}
		cronRepo := repo.NewCron(database.GetDB())
		crons, err := cronRepo.Crons()
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(crons) == 0 {
			fmt.Println("No crons found")
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintf(w, "ID\t| Name\t| Run Every\t| Last Run\t| Took Time\n")
		for _, cron := range crons {
			fmt.Fprintf(w, "%d\t| %s\t| %d\t| %s\t| %f\n",
				cron.Id,
				cron.Title,
				cron.RunEvery,
				time.Unix(cron.LastRun, 0).Format(time.RFC3339),
				cron.RunTime,
			)
		}
		w.Flush()
	},
}
