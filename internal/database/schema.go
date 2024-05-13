package database

import (
	"fmt"
	"log/slog"
)

type migrateFn func(*DB) error

var dbMigrate = [...]migrateFn{v1}

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
			tx.Rollback()
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
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tbl_schema_ver (
		version INT(5) UNSIGNED NOT NULL
	)`)
	if err != nil {
		return -1
	}
	version := 0
	db.Get(&version, `SELECT version FROM tbl_schema_ver`)
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
	// table: tbl_admin
	slog.Debug("[DB] Creating table: tbl_admin")
	_, err := db.Exec(`CREATE TABLE tbl_admin (
		id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
		username VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		lastactive_ts INT(11) UNSIGNED NOT NULL DEFAULT 0,
		created_ts INT(11) UNSIGNED NOT NULL DEFAULT 0,
		PRIMARY KEY (id)
	)`)
	if err != nil {
		return err
	}
	slog.Debug("[DB] Adding unique key to table: tbl_admin")
	_, err = db.Exec(`ALTER TABLE tbl_admin ADD UNIQUE KEY (username)`)
	if err != nil {
		return err
	}

	// table: tbl_cron
	slog.Debug("[DB] Creating table: tbl_cron")
	_, err = db.Exec(`CREATE TABLE tbl_cron (
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
	_, err = db.Exec(`INSERT INTO tbl_cron (title, slug, description, run_every)
		VALUES ('Delete old probe logs', 'delete_old_probe_logs', 'Delete old probe log from the database',120);`)
	if err != nil {
		return err
	}

	// table: tbl_node
	slog.Debug("[DB] Creating table: tbl_node")
	_, err = db.Exec(`CREATE TABLE tbl_node (
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

	// table: tbl_prober
	slog.Debug("[DB] Creating table: tbl_prober")
	_, err = db.Exec(`CREATE TABLE tbl_prober (
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
	_, err = db.Exec(`CREATE TABLE tbl_probe_log (
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
