# GitHub Username Checker

Monitor GitHub username availability with optional auto-claim and Telegram notifications.

## Features

- Concurrent username checking via worker pool
- Proxy support (HTTP/SOCKS5)
- Telegram notifications for available usernames
- Auto-claim username via GitHub session cookies

## Setup

1. Copy environment template:
   ```bash
   cp env_example.env .env
   ```

2. Configure `.env` with your GitHub cookies and Telegram bot token

3. Add usernames to check (one per line):
   ```bash
   echo "desired-username" >> data/users_set/usernames.txt
   ```

4. Add proxies to `data/proxy.txt` (optional, format: `host:port`, `host:port:user:pass` or `user:pass@host:port`)

## Usage

```bash
# Run directly
go run .

# Build and run
make build
./build/main

# Build for all platforms
make build-all
```

## Configuration

| Variable | Description |
|----------|-------------|
| `MAX_GOROUTINES` | Number of concurrent workers |
| `RETRY_COUNT` | Request retry attempts |
| `ENABLE_USERNAME_CHANGE` | Auto-claim found usernames |
| `ENABLE_TELEGRAM_MESSAGE_IF_404` | Notify when 404 found |
| `ENABLE_PROXY_CHECK` | Filter dead proxies on startup |

## License

MIT
