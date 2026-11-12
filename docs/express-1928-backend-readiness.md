# Express 1928 — backend readiness review

Date: 2026-11-12. Scope: single-operator Go backend only; **not** a certification of a live event or a finished host interface.

## CI verification

GitHub Actions run `43173463370` for commit `873ed48029b1b0c323ed2a3ab6c040f1956ab9c6` completed successfully. Its workflow executes Go unit/integration tests, race detector and static checks. This resolves the previously failing test fixture path and incorrect `/deduction` URL (`/deduce` is registered).

A separate follow-up commit adds registry configuration tests. Its CI must be checked independently before claiming that the newest revision is green.

## Covered backend capabilities

- Registered `case_express_1928` scenario and cross-file integrity validation.
- Scene navigation, evidence-gated choices, branching conversations and trust.
- Deductions, motives, false Boris accusation and three verdict outcomes.
- Timed events with default actions, story clock and real-time investigation/accusation deadlines.
- JSON autosave, reload and restart persistence, including conversation, timed event and final verdict.
- Single-operator `GET /state` and command endpoints; discovered evidence only in public state.
- Negative tests for invalid actions, corrupt saves and invalid API commands.

## Known gaps before a live game

1. No approved production host interface; backend-only endpoints are not a usable live-show control panel.
2. No documented full in-person rehearsal with 6–15 participants and a shared TV.
3. Story pacing, readability, clues and character dialogue need human playtesting, not just automated reachability checks.
4. No production authentication, TLS, multi-operator synchronization or distributed session storage; run on a trusted machine/network.
5. Legacy API endpoints and error formats are not uniformly standardized.
6. Verify all mandatory clue chains and timed events with the final host UI before announcing live readiness.

## Recommended next milestone

Freeze the backend API after a manual operator rehearsal. Then, **only with explicit approval**, scope the host interface and TV display as separate deliverables. Keep the canonical crime mechanics in `docs/case-43-crime-mechanics-revision.md`.

See `docs/express-1928-operator-runbook.md` and `docs/backend-api.md` for startup and API details.
