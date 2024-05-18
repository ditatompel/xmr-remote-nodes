package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"xmr-remote-nodes/cmd"
	"xmr-remote-nodes/internal/database"
	"xmr-remote-nodes/internal/repo"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Create Admin",
	Long:  `Create an admin account for WebUI access.`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Usage: xmr-nodes admin create")
			os.Exit(1)
		}

		if args[0] == "create" {
			if err := database.ConnectDB(); err != nil {
				panic(err)
			}
			if err := createAdmin(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Admin account created")
			os.Exit(0)
		}
	},
}

func init() {
	cmd.Root.AddCommand(serveCmd)
	cmd.Root.AddCommand(importCmd)
	cmd.Root.AddCommand(probersCmd)
	probersCmd.AddCommand(listProbersCmd)
	probersCmd.AddCommand(addProbersCmd)
	probersCmd.AddCommand(deleteProbersCmd)
	listProbersCmd.Flags().StringP("sort-by", "s", "last_submit_ts", "Sort by column name, can be id or last_submit_ts")
	listProbersCmd.Flags().StringP("sort-dir", "d", "desc", "Sort direction, can be asc or desc")
}

func createAdmin() error {
	admin := repo.NewAdminRepo(database.GetDB())
	a := repo.Admin{
		Username: stringPrompt("Username:"),
		Password: passPrompt("Password:"),
	}
	_, err := admin.CreateAdmin(&a)
	return err
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

func passPrompt(label string) string {
	var s string
	for {
		fmt.Fprint(os.Stderr, label+" ")
		b, _ := term.ReadPassword(int(syscall.Stdin))
		s = string(b)
		if s != "" {
			break
		}
	}
	fmt.Println()
	return s
}
