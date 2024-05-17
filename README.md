# XMR Remote Nodes

Source code of [https://xmr.ditatompel.com](https://xmr.ditatompel.com), a website that helps you monitor your favourite Monero remote nodes.

## Requirements

### Server & Prober requirements

- Go >= 1.22
- Linux Machine (AMD64 or ARM64)

### Server requirements

- MySQL/MariaDB
- [GeoIP Database](https://dev.maxmind.com/geoip/geoip2/geolite2/) (optional). Place it to `./assets/geoip`, see [./internal/repo/geoip.go](./internal/repo/geoip.go).

## Installation

### For initial server setup:

1. Download [GeoIP Database](https://dev.maxmind.com/geoip/geoip2/geolite2/) and place it to `./assets/geoip`. (see [./internal/repo/geoip.go](./internal/repo/geoip.go)).
2. Copy `.env.example` to `.env` and edit it to match with server environment.
3. Build the binary with `make build`.
4. Run the service with `./bin/xmr-nodes-static-linux-<YOUR_CPU_ARCH> serve`.

To create admin user (for creating prober API key from Web-UI, execute `./bin/xmr-nodes-static-linux-<YOUR_CPU_ARCH> admin create`).

Systemd example: [./tools/resources/init/xmr-nodes-server.service](./tools/resources/init/xmr-nodes-server.service).

### For initial prober setup:

1. Create API key for prober
2. Copy `.env.example` to `.env` and edit it to match with prober environment.
3. Build the binary with `make build`.
4. Run the service with `./bin/xmr-nodes-static-linux-<YOUR_CPU_ARCH> prober`.

Systemd example: [xmr-nodes-prober.service](./tools/resources/init/xmr-nodes-prober.service) and [xmr-nodes-prober.timer](./tools/resources/init/xmr-nodes-prober.timer).

