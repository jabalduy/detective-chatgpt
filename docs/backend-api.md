# Backend API contract (single-operator game)

This document describes the **existing** HTTP endpoints. It is not a frontend design.
The game is operated from one laptop; no multiplayer synchronization is required.

## Session and lifecycle

- Server default port: `8080` (override with `PORT`).
- A browser session is identified by the `sessionID` cookie, created by the server.
  Reuse that cookie across requests and restarts to access the same save.
- `POST /start?caseID=<id>` starts a case. Response includes
  `status`, `title`, `message`, `intro`, `knownInfo`.
- `POST /save` writes the current session to disk.
- `POST /load` restores the saved session.
- Successful game-changing POST commands autosave. Certain GET endpoints
  (`/train`, `/event`, `/examinations`) can advance scheduled game
  mechanics and also autosave.
- `GAME_SAVE_DIR` overrides the default `data/saves` directory.
- `GAME_CASE_DIR` optionally points to the parent directory containing the case folders (for example, `cases/case_express_1928`). Without it, case paths are relative to the server working directory.

## Read-only snapshot

`GET /state` returns a single JSON object for the currently running case.
It also processes due train stops and timed events before returning the snapshot, and autosaves the resulting state. Representative response:

```json
{
  "case_id": "example",
  "started": true,
  "time": "2026-01-01T23:50:00Z",
  "paused": false,
  "scene": {
    "id": "corridor",
    "title": "Corridor",
    "background": "",
    "description": "",
    "choices": [{"id": "inspect", "text": "Inspect", "kind": "action"}]
  },
  "characters": {},
  "found_clue_ids": [],
  "found_fact_ids": [],
  "found_motive_ids": []
}
```

Optional fields include `train`, `conversation`, `active_event`, and
`ending`. The snapshot exposes **available** scene/dialogue/event choices,
not their hidden requirements or consequences. The returned clue/fact/motive
arrays contain IDs, not their full descriptions.

**Important:** `GET /state` is a snapshot, not a command. Use the
corresponding command endpoint to perform an action; do not infer that a
choice succeeded until its response is successful.

## Scenes and conversations

| Method | Endpoint | Purpose |
| --- | --- | --- |
| GET | `/scene` | Current scene and available choices |
| POST | `/scene/choose?choiceID=<id>` | Execute scene choice |
| GET | `/characters` | Current character state |
| POST | `/conversation/start?id=<id>` | Start a branching conversation |
| GET | `/conversation` | Current conversation node and choices |
| POST | `/conversation/choose?choiceID=<id>` | Choose dialogue branch |

A conversation pauses the story clock; completing it resumes the clock.
Scene changes are blocked during an active conversation or event.

## Train, timed events, and investigation

| Method | Endpoint | Purpose |
| --- | --- | --- |
| GET | `/time` | Game clock |
| GET | `/train` | Train schedule/status; advances due station events |
| GET | `/event` | Current event; advances due events; may return 204 |
| POST | `/event/choose?optionID=<id>` | Choose an event response |
| POST | `/event/expire` | Apply an expired event's default consequence |
| GET | `/examinations` | Examination status; completes due examinations |
| POST | `/examinations/start?id=<id>` | Begin an examination |
| GET | `/evidence` | Available evidence presentations |
| POST | `/evidence/present` | Present evidence (see handler for parameters) |
| POST | `/accuse` | Submit accusation JSON |
| GET | `/ending` | Final outcome; may return 204 |

Accusation request JSON uses `suspect_id`, `motive_id`, and `fact_ids`.
The ending is terminal: no further game actions are allowed.

## Other existing routes

The server also registers `/cases`, `/case`, `/locations`,
`/objects`, `/inspect-object`, `/suspects`, `/dialogues`,
`/ask-dialogue`, `/dialogue/start`, `/dialogue/finish`,
`/confront`, `/deduce`, `/facts`, and `/motives`.
Their request/response contracts should be documented separately after
reviewing each handler.

## Error handling and integration cautions

- 400: missing or malformed input on endpoints that validate it.
- 404: missing case or save, where applicable.
- 405: wrong HTTP method.
- 409: game not started, invalid game transition, or blocked choice.
- 422: invalid or unsupported saved game on load.
- 500: server or persistence error.

These are **common statuses, not a universal schema**: some older endpoints
use different error messages/statuses. Error bodies are currently often plain
text rather than structured JSON. The API is not yet frozen.

`GET /state` processes due train stops and timed events, but not examination completion. Consumers should poll `/examinations` for examination progress.
