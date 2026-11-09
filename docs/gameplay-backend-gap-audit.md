# Gameplay-to-backend gap audit

Source of requirements: [game-design-spec.md](game-design-spec.md).
This audit distinguishes **implemented**, **partial**, and **missing**;
it does not authorize frontend or story changes.

| Requirement | Status | Evidence / remaining work |
|---|---|---|
| Single operator, one browser session | Implemented | Session-cookie HTTP server; no player accounts |
| Variable number of participants | Implemented by design | No per-player game state; no fixed player count needed |
| Scenes and branching choices | Implemented | `internal/game/scene.go` |
| Conversations, evidence presentation, trust | Partial | Existing conversation/evidence engine; character-specific question budgets need implementation |
| 60-minute wall-clock investigation deadline | Missing | `internal/game/clock.go` only tracks in-world clock |
| Five-minute accusation deadline | Missing | `/accuse` has no real-time phase gate |
| Allow current conversation to finish at deadline | Missing | Requires phase/deadline state and guards on new actions |
| In-world time costs and train schedule | Implemented | Clock, actions, stations |
| Ten-second urgent choices | Implemented foundation | Event wall-clock deadlines and default outcomes; authoring 3–4 such events remains |
| Progressive carriage unlocking | Partial | Open locations/actions supported; authored unlock routes and map presentation not yet built |
| Interactive illustrated room objects | Partial | Object inspection exists; clickable artwork is frontend work requiring approval |
| Automatic notebook and timeline | Partial | Evidence/facts/character APIs exist; unified testimony and chronology view not established |
| Manual logical deduction (no combine mechanic) | Requires adjustment | Legacy `/deduce` endpoint exists; do not surface as gameplay |
| Unlimited escalating-cost hints | Missing | Need hint definitions, cost progression, API, persistence, tests |
| One final accusation only, after investigation | Partial | Ending blocks repeated accusation, but `/accuse` currently permits early submission |
| Multiple endings | Implemented foundation | Ending engine; authored outcome logic still needed |
| Operator music control | Missing | Frontend/audio feature, deferred pending approval |
| Cinematic final reconstruction | Missing authored assets | Narrative and audiovisual design deferred |

## Proposed engineering sequence

1. Add a **persisted real-time investigation phase** with deadlines and
   conversation grace period. Keep it distinct from in-world train time.
2. Add phase guards to investigation endpoints and tests covering deadlines,
   restarts, event deadlines, and the final accusation.
3. Add configurable question budgets and persistent usage counters.
4. Add hint costs/usage persistence and a spoiler-safe hint API.
5. Provide notebook read models (testimony, chronology, evidence) without
   changing frontend.
6. Verify case-loading constraints and test the complete authored mystery
   once it exists.

**Important:** Wall-clock deadlines must use absolute persisted timestamps
so server restarts cannot grant extra investigation time. The 5-minute
accusation timer should start at the investigation deadline (or when a
grace-period conversation finishes, if that is the chosen rule); the exact
grace-period interaction needs a deterministic engine specification.
