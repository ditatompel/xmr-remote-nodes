package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"xmr-remote-nodes/internal/database"
	"xmr-remote-nodes/internal/repo"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var adminCmd = &cobra.Command{
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
	rootCmd.AddCommand(adminCmd)
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
