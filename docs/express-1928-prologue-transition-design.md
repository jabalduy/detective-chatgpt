# Express 1928 — prologue-to-investigation architecture (stage 84)

**Status:** approved direction for the next backend iteration, **design only**. No Go runtime, JSON case, or frontend changes in this stage. Based on the canonical murder mechanics in `docs/case-43-crime-mechanics-revision.md` and chronology audit stage 83.

## Decision: preserve a playable 21:00 prologue

The case opens at **21:00** while Victor is alive. The murder happens at **21:39** in his compartment, Boris discovers the body privately at **22:03**, the train reaches the tunnel at **22:10**, and the body is officially discovered at **22:17**. The **60-minute real-time investigation countdown starts at the transition after 22:17**, not when the prologue begins. The 5-minute accusation window follows investigation. Prologue duration should be short and operator-paced, with a configurable upper bound determined in rehearsal.

**Important:** This is a future behavior specification. Today `POST /start` starts the 60-minute investigation immediately, and evidence scenes are open at 21:00. The contradiction is **not yet fixed in code**.

## Playable acts

### Act 0: Arrival and prologue (21:00–22:17 in story time)

- Open with carriage reception, Sophie, Marta's intact numbered uniform ribbon, Victor alive and the brass locomotive paperweight visibly on his desk.
- Allow conversations and observations that do not reveal post-murder evidence; characters' positions and testimony reflect **that moment**, not later events.
- 21:12 Max/Klara quarrel and 21:24 Sophie departure can be presented as timed interruptions.
- 21:39 murder happens **offstage**. Do not reveal the killer, the act, or the body to players.
- 21:51 Marta at the service journal is an ambiguous observation, not an accusation.
- 22:09 hooded Graf and optional rooftop pursuit occur before the 22:10 tunnel; event 105 only follows the pursuit choice.
- 22:17 official discovery is a **mandatory story transition**: stop free prologue interaction, play discovery announcement, open the investigation locations and evidence. The transition is not a skippable optional event.

### Act 1: Investigation (starts after 22:17)

- Open Victor's compartment evidence (ribbon fragment, missing paperweight, stopped watch, Sophie note), staff box, Nicolas' route log, Leo documents, interrogation confrontations and deductions.
- Keep all existing mandatory proof paths and optional evidence alternatives; the chase must never be required for conviction.
- Start wall-clock 60-minute investigation deadline **at act transition**; keep a separate in-world route clock, which may continue after 22:17.
- The operator can view what players actually witnessed during prologue; unobserved events do not automatically become known evidence.
- No forced disclosure that Boris altered the watch until the existing confrontation/deduction chain.

### Act 2: Accusation and ending

- Existing 5-minute accusation window, server-authoritative verdict and immutable ending remain unchanged.

## State model proposal

Introduce a distinct **narrative act** independent of `RealTimePhase`:

- `prologue`
- `investigation`
- `finished` (narrative completion; wall-clock accusation still handled by `RealTimePhase`)

Persist the narrative act and a **transition-completed marker** with the save. Never infer act solely from the current scene or client local time. Older saves require explicit migration/default handling; do not silently reinterpret old investigation saves as prologues.

Define a server-owned `BeginInvestigation` transition that is idempotent under session mutex:

1. Verify the story clock has reached the official discovery and that no urgent choice/conversation is being interrupted; if needed, queue the transition.
2. Resolve pending pre-discovery events according to existing timeout/default rules, without granting missed optional evidence.
3. Mark prologue complete, move to discovery/investigation scene, update character availability and reveal only post-discovery scene choices.
4. Set `InvestigationDeadline = now + 60m`, `AccusationDeadline = InvestigationDeadline + 5m` **once**.
5. Autosave the entire transition atomically before acknowledging success.

Do **not** simply shift `startTime.json` to 22:17: that would skip the pre-murder observations and rooftop chase. Do **not** add a 77-minute real-time wait; story time can advance through authored prologue beats.

## Backend integration points

| Existing code/data | Required future change |
| --- | --- |
| `internal/game/realtime_phase.go` | Prologue must not be treated as investigation; add explicit prologue handling or a separate guard for `CanInvestigateAt` |
| `internal/game/state.go` | Narrative act, discovery transition marker, prologue observations |
| `internal/server/start_handler.go` and `newGame` | Start in prologue without arming 60-minute deadline |
| `internal/game/scene.go` | Gate post-murder scene choices on narrative act; current `RequirementsMet` does not support act gating |
| `internal/game/event.go` | Retain 21:12/21:24/21:51/22:09 triggers; make discovery mandatory, not a time-limited event |
| `internal/server/persistence.go` | Save/restore act and transition state; maintain backward compatibility |
| `internal/server/state_handler.go` | Expose player-safe act and current scene; keep secret evidence hidden |
| `cases/case_express_1928/scenes.json` | Author dedicated pre-murder reception/compartment-preview scenes and investigation-only evidence |
| `cases/case_express_1928/conversations.json` | Avoid witnesses describing a murder before 22:17 |
| `docs/backend-api.md` | Document new phase/transition behavior once implemented |

**API caution:** The existing `phase` field is a real-time phase. Frontend must not treat it as the narrative act. Any new field/route must be implemented and documented before use.

## Important edge cases

1. **Clock advancement:** a single operator action may jump over 21:39, 22:09 and 22:17. Process mandatory story milestones in chronological order, not merely the latest clock time; ensure optional rooftop event is not silently skipped if its trigger is crossed.
2. **Conversation pause:** in-world clock pauses during dialogue. The prologue must not deadlock because a conversation never ends; give the operator a non-destructive way to close it.
3. **Active rooftop event at 22:17:** preserve the 10-second server deadline; resolve or queue discovery, never erase the event.
4. **Reload at 22:17:** investigation deadline must not restart; repeated transition attempts are no-ops.
5. **Pre-murder evidence leakage:** the operator must not be able to access fact 1/2/3/5/13 or the body from the prologue scene, REST evidence routes, or deductions. Scene gating alone is insufficient.
6. **Boris timeline:** the 21:24 Sophie interruption cannot be motivated by Boris' later 22:03 discovery; treat it as independent pre-existing protectiveness or rewrite it.
7. **Story pacing:** 21:00 to 22:17 is **77 minutes of story time**, not 77 minutes of wall time. Rehearse a condensed prologue; document how route time advances without confusing players.
8. **Safety/realism:** roof access must be validated against the chosen carriage prototype or adapted to an exterior platform.

## Proposed incremental implementation sequence

1. Add narrative act to game state and persistence, with backward-compatible tests.
2. Add server-authoritative prologue transition and deadline initialization, with duplicate/concurrent request tests.
3. Create prologue-only scenes, character dialogue, and act-gated evidence/actions across **all** routes.
4. Add mandatory 22:17 discovery beat, ensure 22:09 chase can resolve before it, and run story-clock skip tests.
5. Run full-case win/unsolved/fail playthroughs through the prologue, plus save/load and real-time deadline tests.
6. Conduct a 6–15-person rehearsal before finalizing pacing or UI.

## Acceptance criteria

- Victor is alive and no murder-only evidence is available at 21:00.
- The murder at 21:39 remains offstage; no player-visible spoiler.
- The 22:09 chase remains optional, resolves before/around tunnel entry, and cannot strand the game.
- At 22:17 the official discovery reliably opens investigation and starts the 60-minute wall-clock timer exactly once.
- Repeated requests and save/restore never reset deadlines or reveal hidden evidence.
- All previously valid evidence routes and verdicts remain possible without the chase.

**Next stage:** implement only the minimal narrative-act state and persistence foundation, with tests, before touching scenes or frontend.
