# linkding-telegram

> Disclaimer: This documentation was generated with the assistance of GPT-5.5.

A small Telegram polling worker that saves links from Telegram messages to [linkding](https://github.com/sissbruecker/linkding).

The application does **not** expose an HTTP server and does not use Telegram webhooks. It continuously polls Telegram Bot API updates, filters messages by optional chat IDs, extracts URLs from messages, and creates bookmarks in linkding.

## How it works

```text
Telegram getUpdates polling -> optional chat ID filter -> URL extraction -> Linkding bookmark API
```

URLs are extracted from:

- Telegram `link_preview_options.url`
- message `entities`
- message `caption_entities`
- plain message `text`
- plain message `caption`

If creating a bookmark in linkding fails, the worker exits with a fatal error so Docker Compose can restart or surface the failure depending on your runtime settings.

## Requirements

- Docker
- Docker Compose
- A Telegram bot token from [@BotFather](https://t.me/BotFather)
- A linkding API token

## Configuration

The application is configured with environment variables. When using the provided `docker-compose.yml`, create a `.env` file in the project root.

### Required variables

| Variable | Description | Example |
| --- | --- | --- |
| `BOT_TOKEN` | Telegram bot token from BotFather | `123456:ABCDEF...` |
| `LINKDING_ADDRESS` | Base URL of the linkding instance as seen from the worker container | `http://linkding:9090` |
| `LINKDING_USER_TOKEN` | linkding user API token | `abc123...` |

### Optional variables

| Variable | Default | Description |
| --- | --- | --- |
| `BOT_UPDATES_BUFFER_SIZE` | `1` | Internal buffered channel size for Telegram updates |
| `BOT_PERMITTED_CHAT_IDS` | empty | Comma-separated allowlist of Telegram chat IDs. If empty, all chats are accepted |
| `BOT_POLL_INTERVAL_SECOND` | `1` | Delay between Telegram polling requests |
| `LINKDING_DEFAULT_TAG` | empty | Tag added to every created bookmark |
| `LOG_LEVEL` | `info` | Log level: `info`, `warn`, `error`, `fatal`, `panic` |
| `LOG_FILE` | empty | Optional path to write logs to instead of stdout |

## Run with Docker Compose

The repository includes a `docker-compose.yml` that starts two services:

- `linkding-telegram` — this worker
- `linkding` — a local linkding instance exposed on `http://localhost:9090`

> Important: `linkding-telegram` has no published ports because it is a polling worker, not an HTTP server.

### 1. Create a Telegram bot

1. Open [@BotFather](https://t.me/BotFather) in Telegram.
2. Run `/newbot` and follow the instructions.
3. Copy the generated bot token.

### 2. Build the worker image

The compose file expects a local image named `linkding-telegram:latest`.

```sh
docker build -t linkding-telegram:latest .
```

Alternatively, use the provided Make target:

```sh
make build
```

### 3. Start linkding first

Start only the linkding service so you can log in and create an API token:

```sh
docker compose up -d linkding
```

Open linkding in your browser:

```text
http://localhost:9090
```

The default credentials from `docker-compose.yml` are:

```text
Username: admin
Password: 123456789
```

For any real deployment, change these values in `docker-compose.yml` before exposing the service.

### 4. Create a linkding API token

In linkding:

1. Log in as `admin`.
2. Open the user/settings area.
3. Create or copy an API token for the user.
4. Save it for the next step as `LINKDING_USER_TOKEN`.

### 5. Create `.env`

Create `.env` in the project root:

```env
BOT_TOKEN=123456:replace-with-your-telegram-bot-token
LINKDING_ADDRESS=http://linkding:9090
LINKDING_USER_TOKEN=replace-with-your-linkding-api-token

# Optional
LINKDING_DEFAULT_TAG=telegram
LOG_LEVEL=info
BOT_POLL_INTERVAL_SECOND=1
# BOT_PERMITTED_CHAT_IDS=123456789,-1001234567890
```

Why `LINKDING_ADDRESS=http://linkding:9090`?

Inside Docker Compose, containers communicate by service name. The worker container reaches linkding via the Compose service name `linkding`, not via `localhost`.

### 6. Start the full stack

```sh
docker compose up -d
```

Or rebuild and start with Make:

```sh
make run
```

### 7. Check logs

```sh
docker compose logs -f linkding-telegram
```

You should see the worker polling Telegram. Send a message containing a URL to your bot. The URL should appear as a bookmark in linkding.

## Common Docker Compose commands

### Stop services

```sh
docker compose stop
```

### Stop and remove containers

```sh
docker compose down
```

### Rebuild after code changes

```sh
docker build -t linkding-telegram:latest .
docker compose up -d
```

Or:

```sh
make run
```

### View linkding logs

```sh
docker compose logs -f linkding
```

### View worker logs

```sh
docker compose logs -f linkding-telegram
```

## Using an existing linkding instance

If you already run linkding elsewhere, you can remove or ignore the `linkding` service from `docker-compose.yml` and point the worker to your existing instance:

```env
LINKDING_ADDRESS=https://linkding.example.com
LINKDING_USER_TOKEN=your-existing-linkding-token
```

Make sure the `linkding-telegram` container can reach that address.

## Restricting allowed Telegram chats

By default, the bot accepts messages from any chat that can contact it. To restrict this, set `BOT_PERMITTED_CHAT_IDS`:

```env
BOT_PERMITTED_CHAT_IDS=123456789,-1001234567890
```

Use comma-separated Telegram chat IDs. Messages from other chats will be ignored.

## Local development

Run tests:

```sh
go test ./...
```

Run lint:

```sh
make lint
```

Build Docker image:

```sh
make build
```

Run with Docker Compose:

```sh
make run
```

## Notes

- The app uses Telegram long polling (`getUpdates`), not webhooks.
- The app does not expose an HTTP endpoint.
- If linkding is unavailable or returns an error while creating a bookmark, the worker exits with a fatal log message.
- The included linkding credentials in `docker-compose.yml` are suitable only for local testing. Change them for real use.
