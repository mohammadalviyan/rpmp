# AI Orchestration Architecture Decision Record

**Status**: Proposed
**Date**: 2026-09-06
**Authors**: Development Team
**Version**: 1.0

## Document Information

This ADR records how LLM agents should work in the RPMP repo. Goal: the same process on Claude Code and Cursor, without a 16-phase platform or a shared TypeScript package.

Shape adapted from a prior internal orchestration ADR (decisions 1–7 scaled down). Content is for RPMP (Go + Next, docs-first).

Informed by [Deterministic AI Orchestration](https://www.praetorian.com/blog/deterministic-ai-orchestration-a-platform-architecture-for-autonomous-development/). We take gateways, thin agents, and a state schema. We do not run 16 phases literally.

---

## 1. Executive Summary

### Problem

Agents skip invariants (`SELECT *`, vendor terms in the UI, scaffolding `frontend/` during docs-only work). Skills that live only in `.claude/` never reach Cursor.

### Goal

1. Knowledge in `.agents/skills/`
2. Tool veneer in `.claude/` and `.cursor/`
3. Optional state in `.data/` (schema committed, runtime gitignored)

### Architecture at a glance

```
.agents/skills/              # Shared knowledge
  orchestrate/
  gateway-backend/
  gateway-frontend/
  grill-me/
  write-a-skill/

.claude/                     # Claude Code veneer
  agents/                    # thin developer + reviewer defs
  skills/                    # symlinks to .agents/skills/
  hooks/                     # DEFERRED

.cursor/
  skills/                    # symlinks + existing Cursor-only skills
  rules/                     # agents-orchestration.mdc

.data/                       # gitignored except SCHEMA.md
  SCHEMA.md                  # committed
  manifest.yaml              # deferred writer
  feedback-loop.json         # reviewers write when the file exists
```

### Key choices (scaled-down 1–7)

1. Short orchestrate skill, not 16 literal phases
2. Two gateways, no shared lane
3. Four thin agents later (2 developers + 2 reviewers). No `manifest-writer` in v1
4. Hooks deferred until `backend/` / `frontend/` and agents are in daily use
5. `.data/` schema now; writers optional
6. Skills sourced in `.agents/`, symlinked for both tools
7. Docs-first, then backend vertical slice, then frontend. Fail-open when hooks are added

---

## 2. Architectural Decisions

### Decision 1: Short orchestration as a skill

#### Context

A 16-phase loop is ceremony while the repo is docs-only. Skipping process entirely recreates drift.

#### Decision

Auto-invoke `.agents/skills/orchestrate/SKILL.md` on **file-changing** work only. Not for read-only questions.

Phases in v1: **1 Setup, 2 Triage, 3 Intent, 4 Gateway, 8 Implement, 13 Test, 16 Done.**

Triage: TRIVIAL / SMALL / MEDIUM / LARGE. MEDIUM runs `grill-me` as a gate before intent. Do not invent the other Praetorian phases until hooks exist.

Until `backend/` and `frontend/` exist, Phase 8 writes `docs/`, ADRs, and skills only. No app scaffold unless the user asks.

#### Rationale

- Auto-invoke beats hoping someone types `/orchestrate`
- A short table matches docs-first RPMP
- grill-me on MEDIUM catches Dashboard vs Use Cases scope mistakes cheaply

#### Alternatives considered

**Literal 16 phases**

- Rejected for v1. Too much ceremony before there is code.

**Slash command only**

- Rejected: same memory failure the skill exists to fix.

#### Consequences

- Orchestrate description must stay "file-changing" so Q&A does not fire it
- `manifest-writer` is optional/deferred; orchestrator may skip manifest I/O in v1

---

### Decision 2: Two domain gateways, no shared lane

#### Context

v1 has no `packages/shared`. Domains are Go API vs Next UI, plus docs that describe them.

#### Decision

- `gateway-backend` for `backend/**` and backend ADRs/API docs
- `gateway-frontend` for `frontend/**` and frontend UI docs

No third gateway. Point at PRD sections; do not paste the whole PRD.

#### Rationale

- Two skill stacks (Go/sqlc vs Next/RSC)
- Shared DTOs are not a package yet; the API JSON is the contract

#### Alternatives considered

**One gateway**

- Rejected: reviewers and invariants differ (SQL vs UI copy)

#### Consequences

- Adding a worker/cron later means a new gateway
- Cross-cutting API+UI work runs backend first, then frontend

---

### Decision 3: Four thin agents, no manifest-writer in v1

#### Context

Need implement vs review in fresh context. Manifest mutation as a fifth agent is overhead while there is no enforcement.

#### Decision

| Agent                | Role                         | v1        |
| -------------------- | ---------------------------- | --------- |
| `backend-developer`  | Go (or backend docs)         | Defined   |
| `frontend-developer` | Next (or frontend docs)      | Defined   |
| `go-reviewer`        | Read-only backend review     | Defined   |
| `frontend-reviewer`  | Read-only frontend review    | Defined   |
| `manifest-writer`    | `.data/manifest.yaml`        | Deferred  |

Agents stay under 150 lines. First action: load the matching gateway. Developers must not spawn. Reviewers write `VERDICT` and, if `.data/feedback-loop.json` exists, update it.

#### Rationale

- Fresh-context review is the useful Praetorian piece
- Deferring manifest-writer avoids a writer nobody calls

#### Alternatives considered

**No reviewer agents until code exists**

- Rejected: definitions are cheap; reviewers can SKIP when there is no `backend/`

#### Consequences

- Orchestrator spawns reviewers after implementation when source exists
- Without hooks, SKIP/PASS is voluntary. That is accepted for v1

---

### Decision 4: Hooks deferred; fail-open later

#### Context

Hooks are the only non-voluntary Claude mechanism. They also lock sessions if they fail closed on bad JSON.

#### Decision

**Do not install hook scripts in this drop.** When `backend/` and agents are real, copy from the prior orchestration repo and adapt paths:

- `agent-first-enforcement`: block main-thread edits to `backend/**` and `frontend/**`
- `post-edit-dirty-bit`: set dirty domains
- `developer-tests-pass`: `go test` / frontend tests on sub-agent exit
- `feedback-loop-stop`: require reviewer + tests PASS

When added, **fail open** on parse errors (log, allow Stop). Do not deadlock the session on a corrupt `.data/feedback-loop.json`.

#### Rationale

- Hooks without folders trap the main thread with nowhere to delegate
- Fail-open is safer for a small team

#### Alternatives considered

**Install hooks now that rewrite to "spawn backend-developer"**

- Rejected: docs-only main thread would be blocked from the only work that exists

#### Consequences

- v1 process is skill-enforced, not kernel-enforced
- ADR should be updated when hooks land

---

### Decision 5: Tool-agnostic state in `.data/`

#### Context

Cursor and Claude both need a contract if we later enforce dirty bits.

#### Decision

Commit `.data/SCHEMA.md`. Gitignore `.data/` except that file.

`manifest.yaml` v1 fields: `current_phase`, `completed_phases`, `current_task`, `findings`, `decisions`, `dirty_files`, `blockers`. Writer deferred.

`feedback-loop.json` v1: `version`, `dirty_domains.backend/frontend`, `domain_status.{backend,frontend}.{reviewer_status,reviewer_notes,tests_status,tests_notes}`.

Reviewers write the JSON when the file exists. They do not create a hook dependency.

#### Rationale

- Schema-as-contract survives tool differences
- Gitignore avoids per-machine merge noise

#### Alternatives considered

**State only in `.claude/`**

- Rejected: Cursor would not see it

#### Consequences

- Schema version field is `1`
- `/resume` not in v1

---

### Decision 6: Skills in `.agents/`, veneer in `.claude/` and `.cursor/`

#### Context

This repo already has Cursor skills (`bro`, `teach`, `unslop`, `technical-writing`). Orchestration skills must not replace them.

#### Decision

Source: `.agents/skills/<name>/SKILL.md`.

Symlink the five orchestration skills into `.claude/skills/` and `.cursor/skills/`:

`orchestrate`, `gateway-backend`, `gateway-frontend`, `grill-me`, `write-a-skill`.

Leave existing Cursor skills in place. Cursor rule `agents-orchestration.mdc` points at `AGENTS.md`.

#### Rationale

- One playbook, two UIs
- Progressive disclosure stays in skill files, not in AGENTS.md

#### Alternatives considered

**Skills only under `.cursor/`**

- Rejected: Claude Code would miss them

#### Consequences

- New skills default to `.agents/` first, then two symlinks

---

### Decision 7: Docs-first, then backend vertical slice

#### Context

No `backend/` or `frontend/` yet. Building hooks and sqlc together would mix product discovery with scaffolding.

#### Decision

Order:

1. This docs/skills drop (current)
2. Phase 0 discovery (PRD §33) as docs/plans
3. **Backend-first vertical slice** after the user asks: `cmd` + one dashboard summary path + adapter stub + tests
4. Frontend Dashboard against that API
5. Then hooks, then optional `manifest-writer`

Do not scaffold app dirs during orchestration setup.

#### Rationale

- Adapter choice (API vs DB) is still open
- Backend invariants (SQL, typed errors) are sharper to test first

#### Alternatives considered

**Generate Next+Go hello-world now**

- Rejected unless the user asks. AGENTS.md invariant 8.

#### Consequences

- Developer agents must refuse to invent `backend/` during docs tasks
- First real feature updates this ADR if layering names change

---

## 3. Out of scope (v1)

| Primitive              | Why deferred                                      |
| ---------------------- | ------------------------------------------------- |
| 16-phase literal loop  | Docs-first; short table is enough                 |
| Claude/Cursor hooks    | No source tree to enforce                         |
| `manifest-writer`      | No Stop hook consuming the manifest               |
| GitNexus / Serena / RTK / beads | Not part of RPMP v1                    |
| Shared TS package      | PRD does not require it                           |
| Compaction / escalation advisor | Solo-scale                                 |

---

## 4. Open questions

1. When to turn on hooks (after first Go PR?)
2. Whether sqlc output path should be a gateway invariant
3. Dual-reviewer race on `feedback-loop.json` (serial spawn in v1)

## 5. References

- `AGENTS.md`
- `.data/SCHEMA.md`
- `docs/PRD.md` §34
- `.agents/skills/orchestrate/SKILL.md`
