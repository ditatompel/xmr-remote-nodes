# XMR Remote Nodes

[![Test](https://github.com/ditatompel/xmr-remote-nodes/actions/workflows/test.yml/badge.svg)](https://github.com/ditatompel/xmr-remote-nodes/actions/workflows/test.yml)
[![BUild](https://github.com/ditatompel/xmr-remote-nodes/actions/workflows/build.yml/badge.svg)](https://github.com/ditatompel/xmr-remote-nodes/actions/workflows/build.yml)
[![Release Binaries](https://github.com/ditatompel/xmr-remote-nodes/actions/workflows/release.yml/badge.svg)](https://github.com/ditatompel/xmr-remote-nodes/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ditatompel/xmr-remote-nodes)](https://goreportcard.com/report/github.com/ditatompel/xmr-remote-nodes)

Source code of [https://xmr.ditatompel.com](https://xmr.ditatompel.com),
a website that helps you monitor your favourite Monero remote nodes.

> :warning: :construction: This project is not mature enough :construction:,
> If you want to use it on your server, please use it with caution.

## How does it work?

Apart from CPU architecture type, you can build two types of binaries from
this project: a **server** and a **client**.

The **clients** is used to fetch node information given by the server. First,
it will ask the server which node to fetch. Then, it will fetch the information
and report back to the server.

The **server** serves the Web UI and the `/api` endpoint that is used by the
clients.

## Requirements

To build the executable binaries, you need:

- Go >= 1.25
- Bun >= 1.1.26
- [a-h/templ][templ-repo] v0.3.977

> **Note**:
>
> - If you want to contribute to the code, please use exact templ version
>   (v0.3.977).

### Server & Prober requirements

- Linux Machines (AMD64 or ARM64)

### Server requirements

- MySQL/MariaDB
- [GeoIP Database][geoip-doc] (optional). Place it to `./assets/geoip`,
  see [./internal/ip/geo/geoip.go](./internal/ip/geo/geoip.go).

## Installation

### For initial server setup:

1. Download [GeoIP Database][geoip-doc] and place it to `./assets/geoip`.
   (see [./internal/ip/geo/geoip.go](./internal/ip/geo/geoip.go)).
2. Pepare your MySQL/MariaDB.
3. Copy `.env.example` to `.env` and edit it to match with server environment.
4. Build the binary with `make server` (or `make build` to build both
   **server** and **client** binaries).
5. Run the service with `./bin/xmr-nodes-server-linux-<YOUR_CPU_ARCH> serve`.

Systemd example: [xmr-nodes-server.service][server-systemd-service].

### For initial prober setup:

1. Create API key for prober
2. Copy `.env.example` to `.env` and edit it to match with prober environment.
3. Build the binary with `make client` (or `make build` to build both
   **server** and **client** binaries).
4. Run the service with `./bin/xmr-nodes-client-linux-<YOUR_CPU_ARCH> probe`.

Systemd example: [xmr-nodes-prober.service][prober-systemd-service] and
[xmr-nodes-prober.timer][prober-systemd-timer].

## Development and Deployment

1. Clone or fork this repository.
2. Prepare the assets: `make prepare`,
3. Run `air serve` (live reload using [air-verse/air][air-repo]).

See the [Makefile](./Makefile).

## ToDo's

- :white_check_mark: Accept IPv6 nodes.
- :white_check_mark: Use `a-h/templ` and `HTMX` instead of `Svelte`.
- :white_check_mark: Accept I2P nodes.
- :white_check_mark: Support Tor hidden service (beta, inform via HTTP header).
- :white_check_mark: Integrate with [Rucknium's Monero Network Scan][rucknium-network-scan]
- :white_check_mark: Implement ban by IP address using [Boog900's Monero Ban List][boog900-ban-list]
- Use Go standard `net/http` instead of `fiber`.

## Acknowledgement

The creators and contributors of these projects have provided valuable
resources, which I am grateful for:

- [jtgrassie/monero-pool][jtgrassie-monero-pool]
- [Rucknium/xmrnetscan][rucknium-network-scan-repo]
- [Boog900/monero-ban-list][boog900-ban-list]
- [Boog900/p2p-proxy-checker][boog900-p2p-proxy-checker]
- [rclone/rclone][rclone]

## Similar Projects

- [lalanza808/monero.fail][monerofail-repo]
- [cake-tech/upptime-monerocom][uptime-monerocom-repo]

## Donation

The servers costs are currently covered by myself. If you find this project
useful, please consider making a donation to help cover the ongoing expenses.
Your contribution will go towards ensuring the continued availability of the
website and **my** `stagenet` and `testnet` public remote nodes.

XMR Donation address:

```plain
82ez97EPriX6tnkL5fSv7u93aZV2gHWeMUc9hS3a8Dd6NxGWPZWhruGTFnNyJ4z1yWQMHgtA81i9h9vq7HWoKZLk9EzWryD
```

![](./internal/handler/views/assets/img/monerotip.png)

Thank you!

### View Only Wallet Details

Primary Address:

```
48D1TAcsq4vMoRu8CNaDsPgvUR7SQ1B6TCXfBsJXAbPb7TjeNZbMhAMKatSA64XRTzFckyMcaF19U8Ejr9wy8iwnTCAUU8x
```

View Key:

```
f26d16b3ab8265a8ddd9db99f69cafeee9081555e67eb03a34778f0391962803
```

Restore height:

```
3457614 (2025-07-18)
```

> View only old wallet address can be found in [PR #201](https://github.com/ditatompel/xmr-remote-nodes/pull/201).

## License

This project is licensed under [BSD-3-Clause](./LICENSE) license.

[templ-repo]: https://github.com/a-h/templ "a-h/templ GitHub repository"
[geoip-doc]: https://dev.maxmind.com/geoip/geolite2-free-geolocation-data/ "GeoLite2 Free documentation"
[server-systemd-service]: ./deployment/init/xmr-nodes-server.service "systemd service example for server"
[prober-systemd-service]: ./deployment/init/xmr-nodes-prober.service "systemd service example for prober"
[prober-systemd-timer]: ./deployment/init/xmr-nodes-prober.timer "systemd timer example for prober"
[air-repo]: https://github.com/air-verse/air "Air - Live reload for Go apps"
[jtgrassie-monero-pool]: https://github.com/jtgrassie/monero-pool "A Monero mining pool server written in C"
[rclone]: https://github.com/rclone/rclone "rclone GitHub repository"
[monerofail-repo]: https://github.com/lalanza808/monero.fail "Lalanza808's monero.fail GitHub repository"
[uptime-monerocom-repo]: https://github.com/cake-tech/upptime-monerocom "monero.com uptime GitHub repository"
[rucknium-network-scan]: https://moneronet.info/ "Rucknium's Monero Network Scan"
[rucknium-network-scan-repo]: https://github.com/Rucknium/xmrnetscan "Rucknium's Monero Network Scan Source Code"
[boog900-ban-list]: https://github.com/Boog900/monero-ban-list "Boog900's Monero Ban List"
[boog900-p2p-proxy-checker]: https://github.com/Boog900/p2p-proxy-checker "Boog900's P2P Proxy Checker"
