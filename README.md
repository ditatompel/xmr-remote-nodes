# XMR Nodes

## Requirements

- [GeoIP Database](https://dev.maxmind.com/geoip/geoip2/geolite2/) (place it to `./assets/geoip`, see [./internal/repo/geoip.go](./internal/repo/geoip.go)).

## Installation

For initial server setup:

1. Create database structure and import `tbl_cron` data from [./tools/resources/database](./tools/resources/database).
2. Download [GeoIP Database](https://dev.maxmind.com/geoip/geoip2/geolite2/) and place it to `./assets/geoip`. (see [./internal/repo/geoip.go](./internal/repo/geoip.go)).
3. Copy `.env.example` to `.env` and edit it to match with server environment.
4. Build the binary with `make build`.
5. Run the service with `./bin/xmr-nodes-static-linux-<YOUR_CPU_ARCH> serve`.

Systemd example: [./tools/resources/init/xmr-nodes-server.service](./tools/resources/init/xmr-nodes-server.service).
