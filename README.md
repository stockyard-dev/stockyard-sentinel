# Stockyard Sentinel

**On-call schedule and incident log — who's on call, what incidents happened, what was the resolution**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted developer tools.

## Quick Start

```bash
docker run -p 9150:9150 -v sentinel_data:/data ghcr.io/stockyard-dev/stockyard-sentinel
```

Or with docker-compose:

```bash
docker-compose up -d
```

Open `http://localhost:9150` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9150` | HTTP port |
| `DATA_DIR` | `./data` | SQLite database directory |
| `SENTINEL_LICENSE_KEY` | *(empty)* | Pro license key |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 5 members, 3 incidents/mo | Unlimited members and incidents |
| Price | Free | $4.99/mo |

Get a Pro license at [stockyard.dev/tools/](https://stockyard.dev/tools/).

## Category

Operations & Teams

## License

Apache 2.0
