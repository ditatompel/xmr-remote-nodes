package monero

import (
	"bufio"
	"fmt"
	"log/slog"
	"net/http"
)

// Fetch and store IP addresses from Boog900's ban list to local db
func (r *moneroRepo) FetchBoog900BanList() error {
	resp, err := http.Get("https://raw.githubusercontent.com/Boog900/monero-ban-list/main/ban_list.txt")
	if err != nil {
		slog.Error(fmt.Sprintf("[MRL] Failed to download Boog900's ban list: %s", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("[MRL] HTTP request return with status code:  %d ", resp.StatusCode)
	}

	// turncate tbl_ban_list table
	if _, err := r.db.Exec("TRUNCATE TABLE tbl_ban_list"); err != nil {
		return err
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		ip := scanner.Text()
		_, err := r.db.Exec(`INSERT INTO tbl_ban_list (ip_addr) VALUES (?)`, ip)
		if err != nil {
			slog.Error(fmt.Sprintf("[MRL] Failed to insert ip: %s", err))
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
