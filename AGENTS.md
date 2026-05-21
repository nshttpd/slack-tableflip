# AGENTS.md

This repository contains a simple Go Slack slash command bot for table-flipping text.

## Project Type
- Go application (single `main.go`)
- Uses only the standard library
- CLI flags for configuration; produces a static binary via `go build`

## Essential Commands
- Build binary: `go build -o slack-tableflip .`
- Run binary: `./slack-tableflip --token=<slack-token> --webhook=<incoming-webhook-url>`
- Run directly: `go run main.go --token=... --webhook=...`
- Optional flags:
  - `--port=3000` (default 3000)
  - `--norage` (disable :rage1: emoji in payload)
- Initialize module (already done): `go mod init ...`

## Architecture & Control Flow
- `flag` package parses args (token, webhook, norage, port required for token/webhook)
- `net/http` server on specified port.
- `GET /healthz` → JSON `{ "wow": "such health" }`
- `POST /tableflip` (Slack slash command):
  - Parses form, validates token
  - Channel: directmessage uses channel_id else "#"+channel_name
  - Text flipped via internal `flipText` (reverse + unicode map) or default '┻━┻'
  - Posts form payload to webhook URL
  - Errors: 405 for non-POST, 400 parse, 401 bad token, 500 webhook fail; 200 OK success

## Dependencies & Style
- Pure stdlib: encoding/json, flag, fmt, io, net/http, os, strings
- No external modules beyond go.mod
- No tests, linting, or CI present
- Binary is self-contained (no runtime deps)
- The main branch for the repository is called `trunk`. Not `main` or `master`
- The container image will be built with the `ko` command
- all github actions will be versioned by SHA and not by a version number

## Deployment
- Local binary or `go run`
- Env vars or flags supported via CLI

## Gotchas
- Token and webhook are required (exits if missing)
- flipText implements basic unicode reversal (matches original 'flip' behavior for common chars)
- No persistent state; all request handling in main
- Direct messages use raw channel_id
- Icon_emoji omitted entirely when --norage (not just deleted)
- Old Node files (server.js, package.json, Procfile) remain for reference but are superseded by Go version
