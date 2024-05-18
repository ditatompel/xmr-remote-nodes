package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"xmr-remote-nodes/internal/database"

	"github.com/spf13/cobra"
)

type importClient struct {
	db *database.DB
}

func newImport(db *database.DB) *importClient {
	return &importClient{db: db}
}

type importData struct {
	Hostname    string `json:"hostname"`
	Port        int    `json:"port"`
	Protocol    string `json:"protocol"`
	IsTor       bool   `json:"is_tor"`
	DateEntered int    `json:"date_entered"`
}

var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import Monero nodes from old API",
	Long: `Import Monero nodes from old API.
This command only available during migration process and will be removed in future versions.`,
	Run: func(_ *cobra.Command, _ []string) {
		if err := database.ConnectDB(); err != nil {
			panic(err)
		}
		req, err := http.NewRequest(http.MethodGet, "https://api.ditatompel.com/monero/remote-node?limit=500", nil)
		if err != nil {
			panic(err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			panic(fmt.Errorf("status code: %d", resp.StatusCode))
		}

		response := struct {
			Data []importData `json:"data"`
		}{}

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			panic(err)
		}

		action := newImport(database.GetDB())

		for _, node := range response.Data {
			if err := action.processData(node); err != nil {
				fmt.Println(err)
			}
		}

		fmt.Println("Total Source Data: ", len(response.Data))

		fmt.Println("Done!")
	},
}

func (i *importClient) processData(node importData) error {
	query := `SELECT id FROM tbl_node WHERE hostname = ? AND port = ? AND protocol = ?`
	row := i.db.QueryRow(query, node.Hostname, node.Port, node.Protocol)
	var id int
	err := row.Scan(&id)
	if err == nil {
		// fmt.Printf("Skipping %s://%s:%d\n", node.Protocol, node.Hostname, node.Port)
		return fmt.Errorf("node already exists")
	}

	// insert
	query = `INSERT INTO tbl_node (hostname, port, protocol, is_tor, nettype, ip_addr, date_entered, last_check_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = i.db.Exec(query, node.Hostname, node.Port, node.Protocol, node.IsTor, "", "", node.DateEntered, "[2,2,2,2,2]")
	if err != nil {
		return err
	}

	fmt.Printf("Imported %s://%s:%d\n", node.Protocol, node.Hostname, node.Port)

	return nil
}
