# Express 1928 — chronology audit (stage 83)

Scope: compare implemented story-clock events with the authoritative crime timeline in `docs/case-43-crime-mechanics-revision.md`. This is a consistency review, not a claim that the full dramatic sequence is production-ready.

## Implemented timed events

| Event ID | Story time | Scene | Notes |
| --- | --- | --- | --- |
| 101 | 21:12 | Max and Klara argue | Optional evidence, 10-second choice |
| 102 | 21:24 | Sophie leaves salon | Optional trust consequences, 10-second choice |
| 104 | 21:51 | Marta handles service journal | Optional evidence, 10-second choice |
| 103 | 22:09 | Hooded Graf heads to roof | Optional chase entry, 10-second choice |
| 105 | Choice-triggered, no fixed `at` | Rooftop/tunnel decision | Optional 10-second follow-up; default is retreat |

The 22:09 chase precedes the canonical 22:10 tunnel. The rooftop follow-up has no independent clock trigger and must only occur if the operator chooses pursuit in event 103. The final evidence chain does not require these events.

## Canonical murder timeline (stage 43)

- 21:15 Sophie leaves invitation note; 21:22 cancels meeting through Nicolas.
- 21:28 Marta overhears Victor and Klara; 21:35 confrontation in small salon.
- 21:38 Marta follows Victor into his compartment; **21:39** she kills him there.
- 21:41–21:46 weapon hidden in service metal-polishing box.
- 22:03 Boris discovers body; **22:06** he alters/stops Victor's watch.
- **22:10** tunnel; **22:17** official discovery.

## Remaining narrative and scheduling gaps

1. **Investigation begins at 21:00, before the canonical murder (21:39) and official discovery (22:17).** The current case lets the operator inspect Victor's body-related evidence immediately. This is a significant story/interaction contradiction, not fixed by event time corrections. A future narrative design decision must either introduce a pre-murder prologue with gated evidence and transition to investigation after 22:17, or explicitly shift the playable starting point and reschedule/represent early events as flashbacks or witness reports. Do not silently move the murder time.
2. **Event 102 at 21:24** describes Boris interrupting Sophie. It is not necessarily contradictory, but the narrative should explain why he intervenes before discovering Victor's body at 22:03; this cannot be motivated by his later cover-up yet.
3. **Event 104 at 21:51** occurs after the murder and concealment, but before the official discovery. Its wording implies an ongoing investigation and should be reviewed once the playable start is resolved.
4. **The roof scene is abstractly authored.** No historical carriage access mechanism has been established. Stage 43 requires verifying a real carriage prototype or moving the chase to an exterior platform.
5. **Story-clock event scheduling is not identical to real-time deadlines.** Long conversations pause the story clock, and events can be discovered late as the route clock catches up. UI should not claim these are guaranteed to occur at exact wall-clock times.

## Automated guard

`internal/cases/express_1928_chronology_test.go` checks all five event triggers, the one-minute separation from tunnel entry, and the choice-triggered rooftop follow-up.

## Recommended next work

Resolve the **21:00 playable start vs 22:17 official discovery** contradiction before locking the operator UI or commissioning final art. Keep Marta as killer, Victor's body in his compartment, Boris's 22:06 watch manipulation, and the 22:10 tunnel unchanged.
