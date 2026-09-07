# Presentation UI port from bot-metrics-center

**Date:** 2026-09-07  
**Status:** Approved for presentation prototype (not MVP backend scope)

## Goal

Ship an RPMP admin UI that looks and behaves like `bot-metrics-center` (all menus, sample content) so stakeholders can click a demo. Keep Next.js. Recolor to BRI + Inter. Swap mock data for the Go API later, one provider method at a time.

## How to execute (copy, do not scratch)

Source repo (read-only): `/Users/alviyan/Project/reference/bot-metrics-center`

**Do this**

- Copy the listed files into `frontend/` on each card.
- Keep component structure, class names, Recharts markup, and mock shapes.
- Rewrite only the framework glue:
  - TanStack `Link` / `createFileRoute` / `useRouterState` → Next `Link` / `page.tsx` / `usePathname`
  - `DashboardHeader` and `AppSidebar` into the RPMP protected layout
- Map IRPAT red CSS variables to `docs/design/bri-tokens.css` (Cakrawala primary, Inter already from FE-03).
- Put all sample numbers behind `frontend/lib/data/` providers. Pages must not import `rpa-data` forever; import the provider.

**Do not do this**

- Do not `git clone` or replace `frontend/` with the Start app.
- Do not change ADR to TanStack Start.
- Do not rewrite charts, tables, or forms from a blank file “to match Next style.”
- Do not call the Go API from presentation pages except where a card explicitly says so (Overview may keep the existing API adapter behind the same provider).
- Do not send real email or write settings to the backend.

## Data mode

```text
RPMP_DATA_MODE=mock   # default for presentation
RPMP_DATA_MODE=api    # later, per-method
```

`mock-provider` starts as a port of `src/lib/rpa-data.ts`.  
`api-provider` wraps existing `/api/v1` routes and grows as BE lands.

UI must show a visible **Preview** or **Demo data** hint when `mock` is on.

## Card order

1. FE-05 shell + provider (blocks the rest)
2. FE-06 Overview
3. FE-07, FE-08, FE-09, FE-10 in parallel after FE-05 (and FE-06 if they share dashboard widgets)
4. FE-11 visual parity vs the live sample

## Ports

- RPMP Next: `:3000`
- RPMP Go API: `:8080` by default
- Sample `bot-metrics-center` also often `:8080`. Stop the sample before running the API, or change one port.
