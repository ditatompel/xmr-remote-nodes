# XMR Remote Nodes

[![Lint](https://github.com/ditatompel/xmr-remote-nodes/actions/workflows/lint.yml/badge.svg)](https://github.com/ditatompel/xmr-remote-nodes/actions/workflows/lint.yml)

Source code of [https://xmr.ditatompel.com](https://xmr.ditatompel.com), a website that helps you monitor your favourite Monero remote nodes.

> :warning: :construction: This project is still freaking 0.x.x :construction:,
> If you want to use it on your server, please use it with caution.

## How does it work?

Apart from CPU architecture type, you can build two types of binaries from this project: a **server** and a **client**.

The **clients** is used to fetch node information given by the server. First, it will ask the server which node to fetch. Then, it will fetch the information and report back to the server.

The **server** serves an embedded Svelte static site for the Web UI. It also serves the `/api` endpoint that is used by the clients and the Web UI itself.

## Requirements

To build the executable binaries, you need:

-   Go >= 1.22
-   NodeJS >= 20

### Server & Prober requirements

-   Linux Machines (AMD64 or ARM64)

### Server requirements

-   MySQL/MariaDB
-   [GeoIP Database](https://dev.maxmind.com/geoip/geoip2/geolite2/) (optional). Place it to `./assets/geoip`, see [./internal/geo/ip.go](./internal/geo/ip.go).

## Installation

### For initial server setup:

1. Download [GeoIP Database](https://dev.maxmind.com/geoip/geoip2/geolite2/) and place it to `./assets/geoip`. (see [./internal/geo/ip.go](./internal/geo/ip.go)).
2. Pepare your MySQL/MariaDB.
3. Copy `.env.example` to `.env` and edit it to match with server environment.
4. Build the binary with `make server` (or `make build` to build both **server** and **client** binaries).
5. Run the service with `./bin/xmr-nodes-server-linux-<YOUR_CPU_ARCH> serve`.

Systemd example: [./deployment/init/xmr-nodes-server.service](./deployment/init/xmr-nodes-server.service).

### For initial prober setup:

1. Create API key for prober
2. Copy `.env.example` to `.env` and edit it to match with prober environment.
3. Build the binary with `make client` (or `make build` to build both **server** and **client** binaries).
4. Run the service with `./bin/xmr-nodes-client-linux-<YOUR_CPU_ARCH> probe`.

Systemd example: [xmr-nodes-prober.service](./deployment/init/xmr-nodes-prober.service) and [xmr-nodes-prober.timer](./deployment/init/xmr-nodes-prober.timer).

## Development and Deployment

See the [Makefile](./Makefile).

## Similar Projects

-   [lalanza808/monero.fail](https://github.com/lalanza808/monero.fail)
-   [cake-tech/upptime-monerocom](https://github.com/cake-tech/upptime-monerocom)

## License

This project is licensed under [GLWTPL](./LICENSE).
