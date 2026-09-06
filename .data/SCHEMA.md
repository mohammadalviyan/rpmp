# `.data/` schema (RPMP v1)

Runtime files under `.data/` are **gitignored** except this document. They are per-machine orchestration state for RPMP, not product data.

`manifest-writer` is **deferred**. Agents may skip creating `manifest.yaml` in v1. Reviewers write `feedback-loop.json` only when that file already exists.

## `manifest.yaml`

Optional narrative log for a single top-level task. Fresh file per prompt if a writer is added later.

```yaml
version: 1
current_phase: 1
completed_phases: []
current_task: ""
findings: []
decisions: []
dirty_files: []
blockers: []
```

| Field              | Type     | Meaning |
| ------------------ | -------- | ------- |
| `version`          | number   | Schema version. v1 = `1` |
| `current_phase`    | number   | Orchestrate phase id (1, 2, 3, 4, 8, 13, or 16 in RPMP v1) |
| `completed_phases` | number[] | Phases finished |
| `current_task`     | string   | Short restatement of the user request |
| `findings`         | string[] | Discovery notes (paths, PRD section refs). Do not paste the full PRD |
| `decisions`        | string[] | Choices (gateway, grill-me skip/run, adapter API vs DB) |
| `dirty_files`      | string[] | Repo-relative paths touched |
| `blockers`         | string[] | What stopped the task |

Writer (when added): only `manifest-writer`. Main thread should not edit this file by hand.

## `feedback-loop.json`

Ephemeral enforcement state. Create it when hooks land, or when a reviewer needs a place to write.

```json
{
  "version": 1,
  "dirty_domains": {
    "backend": false,
    "frontend": false
  },
  "domain_status": {
    "backend": {
      "reviewer_status": "PENDING",
      "reviewer_notes": "",
      "tests_status": "PENDING",
      "tests_notes": ""
    },
    "frontend": {
      "reviewer_status": "PENDING",
      "reviewer_notes": "",
      "tests_status": "PENDING",
      "tests_notes": ""
    }
  }
}
```

| Field | Meaning |
| ----- | ------- |
| `dirty_domains.backend` | `true` if Go / `backend/**` (or backend docs-as-code) changed this session |
| `dirty_domains.frontend` | `true` if Next / `frontend/**` changed |
| `domain_status.*.reviewer_status` | `PENDING` \| `PASS` \| `FAIL` |
| `domain_status.*.reviewer_notes` | One-line reviewer summary (max ~200 chars) |
| `domain_status.*.tests_status` | `PENDING` \| `PASS` \| `FAIL` |
| `domain_status.*.tests_notes` | Test runner summary |

There is no shared-package domain in RPMP v1.

If the file is missing, `go-reviewer` and `frontend-reviewer` still return a `VERDICT` in chat. They must not fail the review because JSON was absent. Hooks, when added, must **fail open** on parse errors (see `docs/adrs/ADR_AI_ORCHESTRATION.md`).
