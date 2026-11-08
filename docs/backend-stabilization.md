# Backend stabilization checklist

The engine is designed for **one operator on one laptop**, not multiplayer.
This checklist records what has been verified and what remains before
declaring a complete playable game.

## Automated coverage

The GitHub Actions workflow runs:
- `go test ./...`
- `go test -race ./...`
- `go vet ./...`

Regression tests cover branching scenes/dialogues, blocked transitions,
timed-event expiration and chaining, investigations and examinations,
evidence presentation, accusation/endings, save/restore, restart recovery,
autosave, and prevention of post-ending station/examination effects.

## API and persistence

- The existing endpoints are documented in `docs/backend-api.md`.
- `GET /state` is a read-only snapshot; it does not advance time-dependent
  mechanics.
- `GET /train`, `GET /event`, and `GET /examinations` may update the
  game state and autosave.
- `POST /load` restores a save, catches up due train/event transitions,
  and persists the resulting state.
- HTTP error bodies are not yet consistently JSON across legacy endpoints.
- Autosave failures are logged **after** a successful HTTP response; clients
  cannot assume a 2xx response proves that the state reached disk.

## Known limitations / next phase

1. Narrative validation and dependency/reachability analyzers exist but
   are not yet enforced as a single case-loading gate.
2. Reachability analysis is optimistic and does not fully model mutually
   exclusive decisions, timing, locations, or all requirement types.
3. The complete authored mystery and its assets have not been validated
   through a human playthrough.
4. There has been no frontend design or implementation approval.
5. A final operator-driven rehearsal with the real case, one laptop, and
   the intended display remains necessary.

Passing CI establishes regression-test health, **not** that every possible
story branch is solvable or that the finished game is ready to run.

## Final CI verification

The last code change in the stabilization phase, commit `f07ce64`,
passed GitHub Actions (unit/integration tests, race detector and vet).
The following documentation commit `f95e0f5` also passed.

**Milestone:** the planned backend foundation and stabilization phase are
complete. This is **not** a claim that the final authored case or frontend
is complete; the limitations listed above remain open for the next phase.
