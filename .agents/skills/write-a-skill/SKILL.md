---
name: write-a-skill
description: Create RPMP agent skills with concise instructions and progressive disclosure. Use when creating or revising a SKILL.md, skill references, or deterministic skill scripts.
---

# Writing skills

## Process

1. Gather the task, triggers, use cases, and whether deterministic scripts are needed.
2. Draft a concise `SKILL.md`.
3. Put detailed or rarely needed material in one-level-deep references.
4. Review triggers, terminology, examples, and RPMP path assumptions.

## Structure

```text
skill-name/
├── SKILL.md
├── REFERENCE.md
├── EXAMPLES.md
└── scripts/
    └── helper
```

Only `SKILL.md` is required.

## Frontmatter

```md
---
name: skill-name
description: What the skill does. Use when concrete triggers apply.
---
```

The description is the only text an agent sees before choosing a skill:

- Maximum 1024 characters.
- State the capability, then specific triggers.
- Use third person.
- Name relevant file types or task terms.

## Progressive disclosure

- Keep `SKILL.md` under 100 lines when practical.
- Split material when it has distinct domains or is rarely needed.
- Link references one level deep.
- Add scripts for deterministic validation, formatting, or repeated transformations.
- Do not copy the entire RPMP PRD into a skill. Cite `FR-x.y`, `NFR-x`, or a persona.

## Review checklist

- [ ] Trigger description is specific.
- [ ] Instructions use RPMP terms and repo-relative paths.
- [ ] No time-sensitive assumptions are presented as facts.
- [ ] Examples are concrete.
- [ ] References are one level deep.
- [ ] Existing `.agents/skills/` remains the source of truth.
