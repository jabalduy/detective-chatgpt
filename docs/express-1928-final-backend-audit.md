# Express 1928 — final backend audit (stage 78)

Date: 2026-11-12. Scope: single-operator Go backend, not a certification that the live game is ready.

## Verified CI baseline

GitHub Actions completed successfully for commits:
- `cf0f8605` — alternate-order investigation and locked-choice tests (stage 74).
- `5d389ff9` — restart boundaries, event deadlines, and partial deductions (stage 75).
- `8152608e` — repeated operator requests and post-verdict guards (stage 76).
- `6f0a7192` — concurrent navigation and dialogue requests (stage 77).

Workflow `.github/workflows/go.yml` runs `go test ./...`, `go test -race ./...`, and `go vet ./...`. The latest checks passed at audit time. This is automated coverage, not a full manual rehearsal.

## Backend capability inventory

| Capability | State | Evidence |
| --- | --- | --- |
| Express 1928 case loading and data integrity | Implemented and tested | `internal/cases/integrity.go`, `internal/cases/express_1928_story_consistency_test.go` |
| Evidence-gated navigation and alternative investigation order | Implemented and tested | `internal/cases/express_1928_alternate_playthrough_test.go` |
| Interrogations, contradictions, and deductions | Implemented and tested | `internal/cases/express_1928_playthrough_test.go`, interrogation tests |
| False accusations and suspect-specific endings | Implemented and tested | `internal/cases/express_1928_all_accusations_test.go`, `internal/cases/express_1928_suspect_endings_test.go` |
| Story events, 10-second defaults, and clock separation | Implemented and tested | `internal/game/event.go`, timing tests |
| Autosave, restart, verdict persistence | Implemented and tested | `internal/server/express_1928_persistence_integration_test.go`, `internal/server/express_1928_restore_boundaries_test.go` |
| Duplicate and simultaneous HTTP operator commands | Tested for representative cases | `internal/server/express_1928_duplicate_requests_test.go`, `internal/server/express_1928_concurrent_requests_test.go` |
| Operator-facing REST endpoints | Implemented, not frozen | `internal/server/server.go`, `docs/backend-api.md` |
| Finished operator UI and TV display | Not implemented/approved | Out of scope |

## Release blockers before a real group session

1. **Operator UI and TV display:** no approved, tested host control panel exists. The REST API alone is not the finished game.
2. **Human rehearsal:** run the entire story with a host and representative group (6–15 people), including urgent choices, clue visibility, pacing, and final accusation.
3. **Narrative acceptance:** confirm the authored dialogue, 1918 backstory, physical clue chain, and final reconstruction are clear to players; automated tests cannot assess drama or fairness.
4. **Network safety:** there is no production authentication/TLS or multi-operator coordination. Keep the server on a trusted local setup; `StartServer` currently binds `:<PORT>`, i.e. all interfaces.
5. **API stabilization:** error responses and legacy endpoints are not fully standardized. Agree on the host UI contract before freezing it.

## Non-blocking for a local, single-operator rehearsal

- Distributed sessions, horizontal scaling, and multi-user accounts are unnecessary for the agreed one-laptop format.
- Polished deployment automation and public hosting are not prerequisites for a controlled local rehearsal.

## Recommended handoff checklist

1. Confirm `go test ./...`, `go test -race ./...`, `go vet ./...` on the intended host.
2. Run from repository root using `docs/express-1928-operator-runbook.md`.
3. Keep `GAME_CASE_DIR` pointed at the parent `cases` directory and `GAME_SAVE_DIR` writable.
4. Agree on operator-only vs TV-safe data exposure; never expose raw case JSON to players.
5. Obtain explicit approval before building frontend or changing the operator workflow.
6. Perform one complete host rehearsal, then fix issues found before declaring live readiness.

Canonical crime mechanics: `docs/case-43-crime-mechanics-revision.md`.
