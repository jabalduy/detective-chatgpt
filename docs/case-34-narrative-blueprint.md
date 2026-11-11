# Case 34 — The Last Carriage (narrative blueprint)

Status: **proposed story architecture**, not an implemented or approved final case. No frontend or existing case JSON is changed.

## Premise and opening

1928, a luxury night train crossing continental Europe. A private reception brings together eight people who all have reasons to resent financier **Edmund Vale**. During the first 10–15 minutes of play the group meets passengers, hears a disagreement over a sealed portfolio, and learns the geography of the train.

Vale is found dead in his locked sleeping compartment shortly after a scheduled tunnel passage. The conductor claims the corridor was watched; a stopped pocket watch and a witness's recollection point to the tunnel as the time of death. The door is locked from inside.

## Ground truth — author-only spoilers

Vale was killed **earlier, in the private dining salon**, during a confrontation with his former business partner **Clara Voss**, who is the fixed murderer. The immediate cause was a heavy blow with a brass desk weight; a subsequent fall and the victim's medication obscure the injury. Clara and an unwitting porter moved the body in a linen service trolley into the sleeping compartment. The apparent locked-room scene was staged using a duplicate service key and a door that self-latches.

The stopped watch was deliberately set, not damaged in the struggle. The tunnel is a false anchor: the train's whistle and temporary dimming of lights helped create a memorable but misleading timeline. Clara's real motive is not simple theft: Vale planned to expose her falsification of a relief-fund ledger while concealing his own role in the fraud. Clara feared prison and the collapse of her family's reputation. Her motive explains the crime but does not excuse it.

**Fair-play constraints:** Do not make a technical trick depend on impossible train geometry. Establish the service corridor, trolley route, self-latching door, key custody, and the dining salon's location before asking players to deduce the staging.

## Principal suspects (8)

| Suspect | Public reason for suspicion | Secret / investigatory value |
|---|---|---|
| Clara Voss — former business partner | Public dispute over the portfolio | **Murderer**; false alibi and altered ledger |
| Julian Vale — victim's nephew | Inheritance and gambling debts | Stole a letter, not the murder weapon |
| Dr. Elise Moreau — physician | Prescribed medication to the victim | Can distinguish medication effects from fatal trauma |
| Henri Beaumont — train conductor | Controls compartment access | Covered up a service-key discrepancy to protect his job |
| Marta Klein — head attendant | Had access to linen trolley | Helped move a covered load without knowing it was a body |
| Sofia Bellini — singer | Victim threatened her contract | Heard an argument in the dining salon |
| Arthur Pembroke — diplomat | Secret financial dealings with Vale | Provides a contradictory station-time account |
| Lev Orlov — journalist | Investigating relief-fund fraud | Holds independent records of the ledger and blackmail |

Each suspect needs at least one truthful statement, one withheld detail, and a plausible reason for silence. Not every secret should be a lie or a crime.

## Three independent paths to the central reveal

1. **Physical path:** dining-salon blood trace beneath a rug → missing brass weight → trolley fiber/transfer trace → staged compartment.
2. **Witness path:** singer hears the earlier argument → attendant recalls unusual trolley movement → conductor's key log contradicts the supposed sealed room.
3. **Timeline path:** station timetable and service logs → physician disputes stopped-watch inference → journalist's correspondence fixes Vale's last confirmed living appearance.

Each path must be capable of establishing **earlier death and scene staging** without requiring every clue from another path. Convicting Clara additionally requires linking her to the earlier confrontation and motive via at least two alternative evidentiary combinations.

## Event and irreversible-choice plan

- Event A: unexpected halt; choose whether to inspect the service corridor or question passengers (investigation opportunity changes).
- Event B: the attendant is about to dispose of laundry; **10-second decision**, default laundry is removed; other traces remain.
- Event C: a passenger attempts to destroy a letter; **10-second decision**, default letter lost; journalist's copy survives.
- Event D: a station stop lets a suspect leave the carriage; **10-second decision**, default suspect departs but recorded testimony and physical clues remain.
- Event E: a confrontation over the portfolio alters trust and access; consequences persist.
- Optional Event F: a second timetable discrepancy is announced; **10-second decision**, default delays access rather than blocking the case.

All event defaults must be unconditional, and every major branch must retain a winning route.

## Suggested pacing (60 + 5 real minutes)

- 00–10: introductions, geography, relationship tensions.
- 10–15: discovery of the body and apparent locked-room story.
- 15–35: parallel investigation paths, interviews, first surprises.
- 35–50: timeline collapses; players can establish body relocation.
- 50–60: motive and culprit evidence, final confrontations.
- 60–65: accusation only, then cinematic reconstruction and one of four endings.

Pacing is a playtest target, not a hard-scripted lock on deductions.

## Endings

- **Full solution:** Clara, motive, and sufficiently corroborated evidence; reconstructs murder and staging.
- **Partial solution:** Clara correctly identified but motive or proof incomplete; uncertainty remains in the public record.
- **Wrong accusation:** another suspect blamed; evidence explains why the false theory was tempting.
- **Unresolved:** no timely accusation; the train reaches its destination.

## Required authoring and QA before implementation

- Create clue/fact/motive IDs and an explicit dependency graph, then verify at least three independent routes to the staging revelation and at least two routes to Clara.
- Author 7–8 suspect conversations with budgets, trust thresholds, schedules, and evidence confrontations.
- Write scene/room objects, location unlocks, events with 10-second deadlines and unconditional defaults, hint tiers, and ending rules.
- Validate narrative references and events; integrate validators into case-loading or CI.
- Run route coverage, missed-clue, timed-default, save/restart, and 60+5-minute end-to-end tests.
- Preserve operator-only single-screen architecture. **No frontend implementation without explicit approval.**
