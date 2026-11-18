# telegram-bot

Telegram bot for some useful automations.

## Configuration

| Variable | Required | Default | Meaning |
| --- | --- | --- | --- |
| `TELEGRAM_BOT_TOKEN` | yes | | Bot token from @BotFather |
| `ELEVENLABS_API_KEY` | yes | | ElevenLabs API key |
| `ELEVENLABS_MODEL_ID` | no | `scribe_v3` | Speech-to-text model |
| `ALLOWED_CHAT_ID` | no | empty (all chats) | Comma-separated chat IDs allowed to use the bot |

## Local run

```sh
cp example.env .env   # fill in the token and key
set -a; . ./.env; set +a
go run .
```

Checks:

```sh
go build ./...
go vet ./...
go test ./...
```

## Docker

```sh
docker compose up -d --build
```

The compose file reads `.env` and builds a distroless image from the Dockerfile.

## Deploy

```sh
cd ansible
cp example.inventory.yml inventory.yml   # fill in host and user
cp vars/example.secrets.yml vars/secrets.yml
ansible-vault encrypt vars/secrets.yml
make run user=<host>
```

`make dry-run user=<host>` runs the same playbook in check mode. The playbook
requires docker on the host, clones `main` into `~/telegram-bot`, renders `.env`
(mode `0600`) from the vault, and runs `docker compose up -d --build`.

## Adding a feature

Create a package under `internal/features/` implementing `feature.Feature`
(`Name`, `Register`), then add it to the `features` slice in `main.go`. Nothing
else changes: the registry, config, and other features stay untouched.
