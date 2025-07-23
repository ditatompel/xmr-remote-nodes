package monero

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
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

// Get list of IP addresses (may contain subnets) from local database
func (r *moneroRepo) banList() ([]string, error) {
	var ips []string
	rows, err := r.db.Query("SELECT ip_addr FROM tbl_ban_list")
	if err != nil {
		return ips, err
	}
	defer rows.Close()

	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			return ips, err
		}
		ips = append(ips, ip)
	}

	if err := rows.Err(); err != nil {
		slog.Warn(fmt.Sprintf("[MRL] Ban list Iteration error: %s", err))
	}

	return ips, err
}

// Check if the given IP address is on the blacklist
//
// TODO: Use `netip.Addr` for ips from net/netip package instead of `net.IP`.
func isBannedIP(banList []string, ips []net.IP) bool {
	var prefixes []netip.Prefix

	for _, entry := range banList {
		// Try parsing as prefix first
		if prefix, err := netip.ParsePrefix(entry); err == nil {
			prefixes = append(prefixes, prefix)
			continue
		}

		if addr, err := netip.ParseAddr(entry); err == nil {
			prefixes = append(prefixes, netip.PrefixFrom(addr, addr.BitLen()))
		}
	}

	for _, ip := range ips {
		// Convert net.IP to netip.Addr
		var parsed netip.Addr
		if ip4 := ip.To4(); ip4 != nil {
			var ip4Arr [4]byte
			copy(ip4Arr[:], ip4)
			parsed = netip.AddrFrom4(ip4Arr)
		} else if ip16 := ip.To16(); ip16 != nil {
			var ip16Arr [16]byte
			copy(ip16Arr[:], ip16)
			parsed = netip.AddrFrom16(ip16Arr)
		} else {
			continue // skip malformed
		}

		for _, prefix := range prefixes {
			if prefix.Contains(parsed) {
				return true
			}
		}
	}

	return false
}
