package database

import (
	"fmt"
	"log/slog"
)

type migrateFn func(*DB) error

var dbMigrate = [...]migrateFn{v1, v2, v3, v4, v5, v6}

func MigrateDb(db *DB) error {
	version := getSchemaVersion(db)
	if version < 0 {
		return fmt.Errorf("[DB] can't get database schema version")
	} else if version == 0 {
		slog.Warn("[DB] No database schema version found, creating schema version 1")
	} else {
		slog.Info(fmt.Sprintf("[DB] Current database schema version: %d", version))
	}

	for ; version < len(dbMigrate); version++ {
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		migrateFn := dbMigrate[version]
		slog.Info(fmt.Sprintf("[DB] Migrating database schema version %d", version+1))

		if err := migrateFn(db); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
		if err := setSchemaVersion(db, version+1); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func getSchemaVersion(db *DB) int {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tbl_schema_ver (
			version INT(5) UNSIGNED NOT NULL
		)`)
	if err != nil {
		return -1
	}
	version := 0
	if err := db.Get(&version, `SELECT version FROM tbl_schema_ver`); err != nil {
		return -1
	}
	return version
}

func setSchemaVersion(db *DB, version int) error {
	_, err := db.Exec(`DELETE FROM tbl_schema_ver`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO tbl_schema_ver (version) VALUES (?)`, version)
	return err
}

func v1(db *DB) error {
	slog.Debug("[DB] Migrating database schema version 1")

	// table: tbl_cron
	slog.Debug("[DB] Creating table: tbl_cron")
	_, err := db.Exec(`
		CREATE TABLE tbl_cron (
			id INT(8) UNSIGNED NOT NULL AUTO_INCREMENT,
			title VARCHAR(255) NOT NULL DEFAULT '',
			slug VARCHAR(255) NOT NULL DEFAULT '',
			description VARCHAR(255) NOT NULL DEFAULT '',
			run_every INT(8) UNSIGNED NOT NULL DEFAULT 60 COMMENT 'in seconds',
			last_run INT(11) UNSIGNED NOT NULL DEFAULT 0,
			next_run INT(11) UNSIGNED NOT NULL DEFAULT 0,
			run_time FLOAT(7,3) UNSIGNED NOT NULL DEFAULT 0.000,
			cron_state TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
			is_enabled TINYINT(1) UNSIGNED NOT NULL DEFAULT 1,
			PRIMARY KEY (id)
		)`)
	if err != nil {
		return err
	}
	slog.Debug("[DB] Adding default cron jobs to table: tbl_cron")
	_, err = db.Exec(`
		INSERT INTO tbl_cron (
			title,
			slug,
			description,
			run_every
		) VALUES (
			'Delete old probe logs',
			'delete_old_probe_logs',
			'Delete old probe log from the database',
			120
		);`)
	if err != nil {
		return err
	}

	// table: tbl_node
	slog.Debug("[DB] Creating table: tbl_node")
	_, err = db.Exec(`
		CREATE TABLE tbl_node (
			id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
			protocol VARCHAR(6) NOT NULL DEFAULT 'http' COMMENT 'http | https',
			hostname VARCHAR(255) NOT NULL,
			port INT(6) UNSIGNED NOT NULL DEFAULT 0,
			is_tor TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
			is_available TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
			nettype VARCHAR(100) NOT NULL COMMENT 'mainnet | stagenet | testnet',
			height BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
			adjusted_time BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
			database_size BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
			difficulty BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
			version VARCHAR(200) NOT NULL DEFAULT '',
			uptime float(5,2) UNSIGNED NOT NULL DEFAULT 0.00,
			estimate_fee INT(9) UNSIGNED NOT NULL DEFAULT 0,
			ip_addr VARCHAR(200) NOT NULL,
			asn INT(9) UNSIGNED NOT NULL DEFAULT 0,
			asn_name VARCHAR(255) NOT NULL DEFAULT '',
			country VARCHAR(100) NOT NULL DEFAULT '',
			country_name VARCHAR(255) NOT NULL DEFAULT '',
			city VARCHAR(255) NOT NULL DEFAULT '',
			lat FLOAT NOT NULL DEFAULT 0 COMMENT 'latitude',
			lon FLOAT NOT NULL DEFAULT 0 COMMENT 'longitude',
			date_entered INT(11) UNSIGNED NOT NULL DEFAULT 0,
			last_checked INT(11) UNSIGNED NOT NULL DEFAULT 0,
			last_check_status TEXT DEFAULT NULL,
			cors_capable TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
			PRIMARY KEY (id)
		)`)
	if err != nil {
		return err
	}

	// NOTE: If you need list of public nodes (for example to seed `tbl_node`
	// data for integration test), you can use `public_nodes` command from
	// `monero-wallet-cli` app. Eg testnet public nodes:
	// echo "public_nodes" | monero-wallet-cli --testnet --wallet-file=wallet --daemon-address=testnet.xmr.ditatompel.com:443 --password-file=pass_file

	// table: tbl_prober
	slog.Debug("[DB] Creating table: tbl_prober")
	_, err = db.Exec(`
		CREATE TABLE tbl_prober (
			id INT(9) UNSIGNED NOT NULL AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			api_key VARCHAR(36) NOT NULL,
			last_submit_ts INT(11) UNSIGNED NOT NULL DEFAULT 0,
			PRIMARY KEY (id)
		)`)
	if err != nil {
		return err
	}

	slog.Debug("[DB] Adding unique key to table: tbl_prober")
	_, err = db.Exec(`ALTER TABLE tbl_prober ADD UNIQUE KEY (api_key)`)
	if err != nil {
		return err
	}

	// table: tbl_probe_log
	slog.Debug("[DB] Creating table: tbl_probe_log")
	_, err = db.Exec(`
		CREATE TABLE tbl_probe_log (
			id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
			node_id INT(11) UNSIGNED NOT NULL DEFAULT 0,
			prober_id INT(9) UNSIGNED NOT NULL DEFAULT 0,
			is_available TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
			height BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
			adjusted_time BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
			database_size BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
			difficulty BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
			estimate_fee INT(9) UNSIGNED NOT NULL DEFAULT 0,
			date_checked INT(11) UNSIGNED NOT NULL DEFAULT 0,
			failed_reason TEXT NOT NULL DEFAULT '',
			fetch_runtime FLOAT(5,2) UNSIGNED NOT NULL DEFAULT 0.00,
			PRIMARY KEY (id)
		)`)
	if err != nil {
		return err
	}
	slog.Debug("[DB] Adding key to table: tbl_probe_log")
	_, err = db.Exec(`ALTER TABLE tbl_probe_log ADD KEY (node_id)`)
	if err != nil {
		return err
	}

	return nil
}

func v2(db *DB) error {
	slog.Debug("[DB] Migrating database schema version 2")

	// table: tbl_fee
	slog.Debug("[DB] Creating table: tbl_fee")
	_, err := db.Exec(`
		CREATE TABLE tbl_fee (
			nettype VARCHAR(100) NOT NULL DEFAULT '',
			estimate_fee INT(9) UNSIGNED NOT NULL DEFAULT 0,
			node_count INT(9) UNSIGNED NOT NULL DEFAULT 0,
			PRIMARY KEY (nettype)
		)`)
	if err != nil {
		return err
	}
	slog.Debug("[DB] Adding default fee to table: tbl_fee")
	_, err = db.Exec(`
		INSERT INTO tbl_fee (
			nettype,
			estimate_fee,
			node_count
		) VALUES (
			'mainnet',
			0,
			0
		), (
			'stagenet',
			0,
			0
		), (
			'testnet',
			0,
			0
		);`)
	if err != nil {
		return err
	}

	slog.Debug("[DB] Adding majority fee cron jobs to table: tbl_cron")
	_, err = db.Exec(`
		INSERT INTO tbl_cron (
			title,
			slug,
			description,
			run_every
		) VALUES (
			'Calculate majority fee',
			'calculate_majority_fee',
			'Calculate majority Monero fee',
			300
		);`)
	if err != nil {
		return err
	}

	return nil
}

func v3(db *DB) error {
	slog.Debug("[DB] Migrating database schema version 3")

	// table: tbl_node
	slog.Debug("[DB] Adding additional columns to tbl_node")
	_, err := db.Exec(`
		ALTER TABLE tbl_node
		ADD COLUMN ipv6_only TINYINT(1) UNSIGNED NOT NULL DEFAULT '0' AFTER cors_capable,
		ADD COLUMN ip_addresses TEXT NOT NULL DEFAULT '' AFTER cors_capable;`)
	if err != nil {
		return err
	}

	return nil
}

func v4(db *DB) error {
	slog.Debug("[DB] Migrating database schema version 4")

	// table: tbl_node
	slog.Debug("[DB] Adding additional columns to tbl_node")
	_, err := db.Exec(`
		ALTER TABLE tbl_node
		ADD COLUMN is_i2p TINYINT(1) UNSIGNED NOT NULL DEFAULT '0' AFTER is_tor;`)
	if err != nil {
		return err
	}

	return nil
}

func v5(db *DB) error {
	slog.Debug("[DB] Migrating database schema version 5")

	// table: tbl_node
	slog.Debug("[DB] Adding additional columns to tbl_node")
	_, err := db.Exec(`
		ALTER TABLE tbl_node
		ADD COLUMN submitter_iphash CHAR(64) NOT NULL DEFAULT ''
		COMMENT 'hashed IP address who submitted the node'
		AFTER date_entered;`)
	if err != nil {
		return err
	}

	return nil
}

func v6(db *DB) error {
	slog.Debug("[DB] Migrating database schema version 6")

	// table: tbl_rucknium_scan
	// Data from rucknium's Monero Network Scan API
	// Adapted from https://github.com/Rucknium/xmrnetscan/blob/c54748b7835f2f7ecf4673ab699950f9fc6afcdf/R/utils_update_database.R#L367-L385
	// Added ID as primary key and only store neccessary columns for this project
	slog.Debug("[DB] Creating table: tbl_rucknium_scan")
	_, err := db.Exec(`
		CREATE TABLE tbl_rucknium_scan (
			id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
			scan_date VARCHAR(100) NOT NULL,
			connected_node_ip VARCHAR(200) NOT NULL,
			is_spy_node TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
			rpc_domain TEXT NOT NULL DEFAULT 'None',
			mrl_ban_list_enabled TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
			dns_ban_list_enabled TINYINT(1) UNSIGNED NOT NULL DEFAULT 0,
			PRIMARY KEY (id)
		)`)
	if err != nil {
		return err
	}

	slog.Debug("[DB] Adding unique key to table: tbl_rucknium_scan")
	_, err = db.Exec(`
		ALTER TABLE tbl_rucknium_scan
		ADD UNIQUE KEY idx_daily (scan_date, connected_node_ip) USING BTREE,
		ADD KEY scan_date (scan_date);`)
	if err != nil {
		return err
	}

	// table: tbl_node
	// Add new columns for Rucknium's MRL data. All submitted node data will be
	// kept and the `is_archived` column will be used as a query parameter.
	slog.Debug("[DB] Adding additional columns to tbl_node")
	_, err = db.Exec(`
		ALTER TABLE tbl_node
		ADD COLUMN dns_ban_list_enabled TINYINT(1) UNSIGNED NOT NULL DEFAULT 0 AFTER ipv6_only,
		ADD COLUMN mrl_ban_list_enabled TINYINT(1) UNSIGNED NOT NULL DEFAULT 0 AFTER ipv6_only,
		ADD COLUMN is_spy_node TINYINT(1) UNSIGNED NOT NULL DEFAULT 0 AFTER ipv6_only,
		ADD COLUMN is_archived TINYINT(1) UNSIGNED NOT NULL DEFAULT 0 AFTER ipv6_only
		;`)
	if err != nil {
		return err
	}

	// Since the ban lists is not (yet) applicable for ipv6, node from i2p and
	// tor network, set all nodes to initial state (2, not applied).
	slog.Debug("[DB] Updating nodes record to initial state of `not applied (2)`")
	_, err = db.Exec(`
		UPDATE tbl_node
		SET
			dns_ban_list_enabled = 2,
			mrl_ban_list_enabled = 2,
			is_spy_node = 2
		;`)
	if err != nil {
		return err
	}

	// table: tbl_ban_list
	// A table for list of banned IP addresses (support subnets). It's
	// recommended to use Boog900's Monero Ban List:
	// https://github.com/Boog900/monero-ban-list
	slog.Debug("[DB] Creating table: tbl_rucknium_scan")
	_, err = db.Exec(`
		CREATE TABLE tbl_ban_list (
			ip_addr VARCHAR(200) NOT NULL COMMENT 'IP address with subnet is supported',
			PRIMARY KEY (ip_addr)
		)`)
	if err != nil {
		return err
	}

	slog.Debug("[DB] Adding cron jobs for MRL and DNS ban list to tbl_cron")
	_, err = db.Exec(`
		INSERT INTO tbl_cron (
			title,
			slug,
			description,
			run_every
		) VALUES (
			'Fetch Rucknium\'s Node Data',
			'fetch_rucknium_node_data',
			'Fetch and store Rucknium\'s individual_node_data to database',
			43200
		), (
			'Check MRL ban list',
			'check_mrl_ban_list',
			'Check monitored node IP addresses with Rucknium\'s node data',
			300
		), (
			'Fetch static MRL ban list',
			'fetch_static_mrl_ban_list',
			'Fetch and store MRL ban list to database',
			172800
		);`)
	if err != nil {
		return err
	}

	return nil
}
