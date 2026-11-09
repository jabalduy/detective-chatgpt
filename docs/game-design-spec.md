# Game design specification — agreed rules

Status: agreed gameplay requirements; implementation status is **not** implied.
The game is a **single-operator, single-screen, cooperative detective mystery**.
No individual player accounts, hidden roles, traitors, or player-specific devices.
Group size is variable (approximately 6, potentially 15); no mechanic may
depend on a fixed participant count.

## Setting and narrative

- Cinematic but rational murder mystery on a luxury overnight train in the
  1920s–1930s, in the tradition of classic detective fiction.
- Begin with introductions to passengers, then the murder.
- Approximately 7–8 principal suspects; one predetermined murderer.
- The apparent circumstances of the crime (time, location, or method) are
  misleading; the reveal must be surprising but logically justified.
- Complex logical investigation, no mandatory ciphers or puzzle minigames.
- Several independent routes to the truth: physical evidence, testimony,
  timelines. Missing one clue must not make victory impossible.
- 3–4 possible outcomes, including full success, partial success, incorrect
  accusation, and unresolved case.
- The full authored mystery is **not yet defined**.

## Interaction model

- All players see the same TV screen, discuss options aloud, and communicate
  their choice to the operator, who clicks on a laptop.
- A schematic train map is used to navigate unlocked carriages.
- Room illustrations contain interactive objects.
- Only currently available actions are shown.
- Initially 2–3 locations are available; more unlock through progress,
  timed events, or alternative conditions, not a single missable clue.
- Interactive objects support inspection and meaningful interactions.
- NPCs move, change mood/trust, and can become unavailable by schedule.
- Interviews have character-specific question limits. Evidence may unlock
  new topics; questions and opportunities may be permanently lost.
- Players reason from the notebook themselves; **no combine-clues-to-deduce
  feature**. Presenting evidence to suspects is a separate mechanic.
- Notebook automatically records discovered clues, testimony, character
  information, and a timeline. It must not reveal undiscovered secrets.

## Time and consequences

- Real-world investigation deadline: **60 minutes**. It continues while
  players discuss and read ordinary dialogue.
- After investigation ends, provide **5 additional real minutes** for the
  one final accusation. If the deadline occurs during a conversation, allow
  that conversation to finish, but prohibit starting new investigation
  actions afterward.
- The train has a separate in-world clock; actions cost in-world minutes,
  and stations/events can be tied to its schedule.
- Approximately 4–6 major unexpected events; 3–4 decisions have a
  **10-second real-world** deadline and an automatic consequence on timeout.
- Wrong choices can cost time, trust, evidence, and access, but multiple
  viable investigative routes must remain.
- Missed important clues can sometimes be recovered through alternative
  routes.
- Hints are unlimited but progressively more expensive in **in-world**
  minutes. Exact cost progression is a balancing parameter, not yet fixed.

## Ending and presentation

- Only one official accusation, at the end. It must identify murderer and
  motive and support the theory with convincing collected evidence.
- The ending reconstructs the true crime cinematically using scenes,
  images, music, and dialogue.
- Soundtrack and sound effects respond automatically to scenes; the
  operator can manually switch music or adjust/mute it.
- No frontend changes are authorized until the interface is separately
  discussed and explicitly approved.

## Implementation follow-up

Before writing the case or frontend, audit the existing Go engine against
these requirements. In particular, verify the 60+5 real-time deadlines,
conversation grace period, question budgets, hint pricing, notebook,
location unlocking, final accusation gating, and audio control.
Keep story authoring and interface design separate from engine work.
