# Express 1928 — frontend handoff specification (stage 79)

**Status: design and integration specification only. No frontend implementation is authorized or included.**

## 1. Product and hardware constraints

- One operator laptop controls the entire game. A shared TV shows a **player-safe presentation**, not a separate player client.
- Audience: 6–15 players discussing decisions together. No personal logins, phones, or multiplayer synchronization.
- Tone: 1928 luxury night train, serious Agatha Christie-style investigation, restrained theatrical transitions, readable Russian typography.
- Backend: existing Go REST API; keep case mechanics and verdict decisions server-side.
- Game structure: 60 minutes of real-time investigation, then 5 minutes for accusation; 10-second urgent choices; separate in-world train clock. An active conversation may receive deadline grace.
- Story authority: `docs/case-43-crime-mechanics-revision.md`. Do not invent a new murderer or move Victor's body.

## 2. Operator screen (private)

Suggested desktop layout, optimized for a single laptop:

1. **Top bar:** case title, investigation/accusation phase, real-time countdown, in-world train time, paused indicator, save/recovery status.
2. **Left navigation:** current scene, available scene choices, location history, train progress.
3. **Center:** active scene text or active conversation with dialogue choices; show one blocking urgent event above other actions.
4. **Right evidence panel:** discovered facts, motives, deductions and accusation preparation. Never reveal unreached facts.
5. **Bottom controls:** start/load/save, presenter visibility toggle, confirmation for irreversible accusation; show HTTP errors without losing current state.

Operator actions should display an in-flight state and disable duplicate submissions until the server responds. A failed command should refresh `GET /state` before allowing another attempt.

## 3. TV presentation (player-safe)

- **Default:** scene title, atmospheric illustration or abstract background, large readable story text, game clock and a limited list of available decisions.
- **Conversation:** character name, spoken dialogue, optional portrait, choices under group discussion; operator-only controls remain hidden.
- **Urgent event:** clear full-screen interruption with countdown and two or more responses. At expiry, show the consequence of the server-selected default, not a locally fabricated result.
- **Evidence reveal:** only already discovered clue/fact descriptions, and only after the operator deliberately chooses to reveal them.
- **Accusation/finale:** suspense transition followed by the server-returned ending; never pre-render hidden murderer or proof.
- **Accessibility:** large text at TV distance, high contrast, reduced-motion option, readable color-independent status indicators, keyboard-operable operator actions.

**Critical privacy boundary:** `GET /state` is designed to omit hidden action effects and unrevealed solution facts, but it includes operator-specific details such as character state, phase/deadlines, evidence IDs, and an ending. Do not mirror raw JSON or the full operator screen onto the TV. Use an explicit frontend allowlist and a separate presentation model.

## 4. Existing API integration

The following routes are implemented, not a frozen contract. See `docs/backend-api.md` and server handlers for authoritative parameters and response shapes.

| UI operation | Existing backend route | Integration note |
| --- | --- | --- |
| Start case | `POST /start?caseID=case_express_1928` | Preserve server-issued `sessionID` cookie |
| Refresh current state | `GET /state` | May advance due train/event logic; not a pure side-effect-free read |
| Scene decision | `POST /scene/choose?choiceID=...` | Only submit currently visible choice IDs |
| Start interrogation | `POST /conversation/start?id=...` | Requires character at current location |
| Choose dialogue | `POST /conversation/choose?choiceID=...` | On end returns `{"ended":true}` |
| Inspect current conversation | `GET /conversation` | Only when active |
| Urgent choice | `POST /event/choose?optionID=...` | Server enforces deadline |
| Expire event | `POST /event/expire` | Prefer server-driven expiry, do not trust local timer |
| Deduction | `POST /deduce` | JSON `{"fact_ids":[1,8]}` |
| Accusation | `POST /accuse` | JSON `{"suspect_id":1,"motive_id":1,"fact_ids":[10,11,12]}` |
| Save / load | `POST /save`, `POST /load` | Same session cookie across restarts |
| Read ending | `GET /ending` | Final verdict is immutable |

`GET /state` includes `case_id`, `started`, `phase`, `investigation_deadline`, `accusation_deadline`, `time`, `paused`, `scene`, `characters`, `conversation`, `active_event`, `ending`, `found_clue_ids`, `found_fact_ids`, and `found_motive_ids`. Optional sections may be absent.

**Important limitation:** found evidence is represented by IDs in the state snapshot, not full narrative descriptions. Frontend needs a separately reviewed way to resolve *only discovered* descriptions through existing evidence endpoints, or a small dedicated server projection. Do not ship the full `cases/case_express_1928/*.json` data to the browser: it contains hidden clues, deduction prerequisites, and the solution.

## 5. State machine and event behavior

- **Not started:** show start or restore.
- **Investigation:** scene actions, conversations, deductions, evidence. Active conversation pauses story clock; real-time deadline remains authoritative.
- **Urgent event:** overlay takes priority and blocks other actions. Server controls the 10-second countdown and default outcome.
- **Accusation:** after investigation deadline, allow the accusation UI; do not assume an early accusation is accepted.
- **Ended:** show ending and disable all game mutations.
- **Connection interrupted:** keep last safe presentation with a disconnected badge; retry state fetch. Do not invent successful commands, reset clocks, or reveal new clues.

Use the backend's absolute deadlines to compute a visual countdown, but resynchronize on every response; client countdown is indicative only. The server is the authority for all state transitions.

## 6. Suggested screen transitions

1. Opening card: `Ночной экспресс: дело Виктора`.
2. Train corridor -> selected compartment/salon/staff/archive/vestibule/luggage/pantry.
3. Dialogue view with optional evidence-gated confrontation.
4. Evidence reveal animation, then return to scene.
5. Ten-second urgent event interrupt (including Graf hooded figure).
6. Accusation desk: suspect + motive + selected evidence.
7. Server-determined verdict and closing card.

The rooftop chase is a narrative requirement from design docs, but the current case event implementation does **not** yet provide a complete interactive rooftop chase. Do not promise a fully playable chase without an additional approved backend/story task.

## 7. Minimum acceptance criteria

- One laptop operates the game; TV shows no hidden clues, solution, operator controls, or raw JSON.
- Operator can start, investigate, interrogate, deduce, handle an urgent event, accuse, see ending, and restore after restart.
- All interactive choices originate from server-visible options; no fabricated state or client-only evidence.
- Ten-second choices expire correctly even if UI freezes or reconnects.
- The accusation screen cannot submit early or submit twice.
- Screens remain usable at TV distance and on the operator laptop.
- A full rehearsal with 6–15 players confirms pacing, clue legibility, and understandable ending.

## 8. Open decisions before implementation

- Approve the operator layout and TV visual direction.
- Choose frontend stack and whether operator + TV use two browser windows or an extended-display desktop layout.
- Decide which discovered evidence descriptions may be shown on TV and when.
- Confirm whether to implement the rooftop chase as a dedicated interaction or keep the existing timed event.
- Define final API changes, especially safe evidence description retrieval and consistent error responses.

**No UI code should be written until these choices are approved.**
