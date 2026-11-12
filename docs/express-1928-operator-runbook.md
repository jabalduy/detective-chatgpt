# Express 1928 backend: operator runbook (stage 64)

This is a **backend-only** runbook. The production operator UI is not implemented or approved.

## Start from the repository root

Requires the Go version declared in `go.mod`.

```sh
go test ./...
go test -race ./...
go vet ./...
PORT=8080 GAME_SAVE_DIR=./data/saves GAME_CASE_DIR=./cases go run .
```

The server binds to `:<PORT>` (all network interfaces); use a trusted local network, or restrict exposure externally. Default port is 8080. Relative static/template paths also assume the repository root is the working directory. `GAME_CASE_DIR` is the parent directory of named case folders, **not** the folder of an individual case.

## Start and drive the case

Use the same session cookie for all requests. The API has no authentication beyond the session cookie; do not expose it publicly without additional protection.

```sh
curl -i -c /tmp/express-cookie.txt -b /tmp/express-cookie.txt \
  -X POST 'http://localhost:8080/start?caseID=case_express_1928'
curl -s -b /tmp/express-cookie.txt 'http://localhost:8080/state'
curl -s -b /tmp/express-cookie.txt 'http://localhost:8080/scene'
curl -s -b /tmp/express-cookie.txt -X POST \
  'http://localhost:8080/scene/choose?choiceID=victor'
curl -s -b /tmp/express-cookie.txt 'http://localhost:8080/state'
```

The server autosaves successful commands and progression GETs. `GET /state` may process scheduled train stops and events; `GET /examinations` progresses examinations. To restore after a restart, reuse the cookie and call `POST /load`. The save file is keyed to the session ID under `GAME_SAVE_DIR`.

## Runtime behavior and diagnostics

- Investigation lasts 60 minutes of real time; accusation follows for five minutes. An already-active dialogue can finish after the investigation deadline, and the server preserves a full accusation window when needed.
- The train's story clock is separate from the real-time deadlines. Conversations pause the story clock; authored action time costs may advance it.
- Timed events have ten-second server-enforced choices and default consequences on expiry.
- `GET /state` exposes discovered evidence IDs and available choices, not hidden solutions or requirements. Do not expose case JSON files to players.
- Common HTTP statuses: 400 bad input, 404 missing case/save, 405 wrong method, 409 blocked transition, 422 invalid save, 500 internal error. Error bodies are not uniformly JSON.
- A corrupt or unsupported save is rejected. Preserve backups of `GAME_SAVE_DIR`; saves are JSON with atomic temporary-file rename.
- Use GitHub Actions test, race detector and static checks before a live session. A successful test run is not a substitute for an operator rehearsal.

## Known limitations

- No production-ready host UI is included in this stage.
- No multiplayer, authentication, TLS, or production deployment configuration is provided.
- File-backed sessions are intended for a single operator, not horizontal scaling.
- API error responses and legacy routes are not yet standardized.

See `docs/backend-api.md` for endpoint contracts and `docs/case-43-crime-mechanics-revision.md` for the authoritative story logic.
