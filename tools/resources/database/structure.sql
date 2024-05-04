
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
DROP TABLE IF EXISTS `tbl_admin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tbl_admin` (
  `id` bigint(30) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(200) NOT NULL,
  `password` varchar(200) NOT NULL,
  `lastactive_ts` int(11) unsigned NOT NULL DEFAULT 0,
  `created_ts` int(11) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `tbl_cron`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tbl_cron` (
  `id` int(8) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(200) NOT NULL DEFAULT '',
  `slug` varchar(200) NOT NULL DEFAULT '',
  `description` varchar(200) DEFAULT NULL,
  `run_every` int(8) unsigned NOT NULL DEFAULT 60 COMMENT 'in seconds',
  `last_run` bigint(20) unsigned DEFAULT NULL,
  `next_run` bigint(20) unsigned DEFAULT NULL,
  `run_time` float(7,3) unsigned NOT NULL DEFAULT 0.000,
  `cron_state` tinyint(1) unsigned NOT NULL DEFAULT 0,
  `is_enabled` int(1) unsigned NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `tbl_node`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tbl_node` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `protocol` varchar(6) NOT NULL DEFAULT 'http' COMMENT 'http | https',
  `hostname` varchar(200) NOT NULL DEFAULT '',
  `port` int(6) unsigned NOT NULL DEFAULT 0,
  `is_tor` tinyint(1) unsigned NOT NULL DEFAULT 0,
  `is_available` tinyint(1) unsigned NOT NULL DEFAULT 0,
  `nettype` varchar(100) NOT NULL COMMENT 'mainnet | stagenet | testnet',
  `height` bigint(20) unsigned NOT NULL DEFAULT 0,
  `adjusted_time` bigint(20) unsigned NOT NULL DEFAULT 0,
  `database_size` bigint(20) unsigned NOT NULL DEFAULT 0,
  `difficulty` bigint(20) unsigned NOT NULL DEFAULT 0,
  `version` varchar(200) NOT NULL DEFAULT '',
  `uptime` float(5,2) unsigned NOT NULL DEFAULT 0.00,
  `estimate_fee` int(9) unsigned NOT NULL DEFAULT 0,
  `ip_addr` varchar(200) NOT NULL,
  `asn` int(9) unsigned NOT NULL DEFAULT 0,
  `asn_name` varchar(200) NOT NULL DEFAULT '',
  `country` varchar(200) NOT NULL DEFAULT '',
  `country_name` varchar(255) NOT NULL DEFAULT '',
  `city` varchar(200) NOT NULL DEFAULT '',
  `lat` float NOT NULL DEFAULT 0 COMMENT 'latitude',
  `lon` float NOT NULL DEFAULT 0 COMMENT 'longitude',
  `date_entered` bigint(20) unsigned NOT NULL DEFAULT 0,
  `last_checked` bigint(20) unsigned NOT NULL DEFAULT 0,
  `last_check_status` text DEFAULT NULL,
  `cors_capable` tinyint(1) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `tbl_probe_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tbl_probe_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `node_id` bigint(20) unsigned NOT NULL DEFAULT 0,
  `prober_id` int(9) unsigned NOT NULL DEFAULT 0,
  `is_available` tinyint(1) unsigned NOT NULL DEFAULT 0,
  `height` bigint(20) unsigned NOT NULL DEFAULT 0,
  `adjusted_time` bigint(20) unsigned NOT NULL DEFAULT 0,
  `database_size` bigint(20) unsigned NOT NULL DEFAULT 0,
  `difficulty` bigint(20) unsigned NOT NULL DEFAULT 0,
  `estimate_fee` int(9) unsigned NOT NULL DEFAULT 0,
  `date_checked` bigint(20) unsigned NOT NULL DEFAULT 0,
  `failed_reason` text NOT NULL DEFAULT '',
  `fetch_runtime` float(5,2) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `node_id` (`node_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `tbl_prober`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tbl_prober` (
  `id` int(9) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `api_key` varchar(36) NOT NULL,
  `last_submit_ts` int(11) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `api_key` (`api_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

