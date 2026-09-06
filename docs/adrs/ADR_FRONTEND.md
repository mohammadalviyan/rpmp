# Frontend Architecture Decision Record

**Status**: Proposed
**Date**: 2026-09-06
**Authors**: Development Team
**Version**: 1.0

## Document Information

This ADR records proposed frontend decisions for RPMP. PRD §34 suggests React / Next.js and is not locked. MVP screens: Dashboard and Use Cases. Copy and API fields use RPMP language, not Orchestrator vendor terms (PRD §17, §39).

`frontend/` does not exist yet. Do not scaffold the Next app unless the user asks.

## 1. Executive Summary

### Technology stack overview

- **Meta-framework**: Next.js App Router (proposed)
- **UI**: React, Tailwind CSS, shadcn/ui (internal admin UI)
- **Data**: React Server Components for dashboard **reads**; client components for filters, tables, and other interactivity
- **Server state on the client**: TanStack Query is **optional**. Do not add it until a client island needs cache, polling, or shared mutation state
- **Auth**: httpOnly cookies set by the Go API. No tokens in JS
- **Charts/tables**: pick when implementing (for example Recharts + a table primitive). Not locked here

### Key choices

1. Next.js App Router, not a TanStack Start clone
2. RSC for first paint of Dashboard/Use Cases lists; client where the PRD needs filters and detail interaction
3. shadcn/ui + Tailwind for a dense internal console
4. Consume only the RPMP `/api/v1` contract
5. No Orchestrator jargon in labels, empty states, or TypeScript types the UI owns

---

## 2. Architectural Decisions

### Decision 1: Next.js App Router

#### Context

The product is a web-only internal console. PRD suggests React / Next.js.

#### Decision

Use the App Router under planned `frontend/app/`. Route groups for `(auth)` vs dashboard shell. File-based routes; do not hand-edit Next generated artifacts.

#### Rationale

- Fits suggested stack and SSR/RSC without adopting another meta-framework
- Internal SSO later can sit on middleware + cookies

#### Alternatives considered

**SPA Vite + React Router**

- Pros: simpler mental model
- Cons: extra auth/SSR work the App Router already covers

**TanStack Start**

- Pros: typed routing
- Cons: cargo-cult of another repo's stack; not in PRD §34

#### Consequences

- Deploy as the approved internal frontend platform
- BFF-style Route Handlers are **not** the Orchestrator client. If used, they only proxy the Go API with cookies

---

### Decision 2: RSC for reads, client for interactivity

#### Context

Dashboard is mostly read models (summary, trends, use-case list). Filters and detail panes need client state.

#### Decision

**Recommended:** fetch dashboard summaries and use-case lists in Server Components (or server loaders) from the Go API using the incoming cookie. Use Client Components for search/filter, pagination controls, and charts that need browser APIs.

Add **TanStack Query** only when client islands must refetch, poll, or share cache. Do not install it "because dashboards always use Query."

#### Rationale

- First paint can be HTML from the server, which helps the <3s dashboard target (PRD §35) if the API is fast
- Avoids duplicating server state in `useState`
- Query is justified for live refresh, not for a static daily snapshot

#### Alternatives considered

**TanStack Query for every read**

- Rejected as default: extra client bundle and cache invalidation for pages that can be RSC

**RSC only, no client**

- Rejected: Use Cases filters and detail UX need client interactivity

#### Consequences

- Keep secrets off `'use client'` modules
- API types should be imported from a thin `frontend/lib/api` that mirrors the Go JSON contract. Do not invent DTO fields the backend does not return. No `packages/shared` in v1; when Go OpenAPI or a checked-in spec exists, generate or copy from that, do not guess

---

### Decision 3: shadcn/ui and Tailwind

#### Context

Internal admin UI: tables, filters, status chips, layout shell.

#### Decision

Use Tailwind + shadcn/ui primitives (Button, Table, Dialog, Select). Customize with classes, not a parallel CSS architecture.

#### Rationale

- Fast to assemble a dense console
- Matches common Next internal tools

#### Alternatives considered

**Enterprise design system only**

- If the org mandates a kit, replace shadcn. Until then, shadcn is the proposed default

#### Consequences

- Do not put business formulas in components. Display numbers the API already computed

---

### Decision 4: Auth cookies, no frontend secrets

#### Context

Internal users. SSO later.

#### Decision

Browser sends cookies to the Go API (`credentials: 'include'` on client fetches). Session established via login (or future SSO redirect). `GET` of current user from RPMP, not from Orchestrator. Never store refresh tokens in `localStorage`.

#### Rationale

- XSS cannot read httpOnly cookies
- SSO-ready: swap IdP behind the same cookie session

#### Alternatives considered

**Bearer token in memory only**

- Weaker for Next RSC (server still needs the cookie or a server-side session)

#### Consequences

- Middleware can gate `app/` dashboard routes on session
- CORS/cookie domain is a backend concern documented with the API

---

### Decision 5: Language in the UI

#### Context

PRD: business users should not need Orchestrator terminology. Example: vendor `JobState = 3` becomes Failed Execution.

#### Decision

Labels, routes, and frontend types use RPMP words: use case, execution, success rate, health, failure. Adapter-only names stay out of `frontend/`.

#### Rationale

- Product positioning (information layer)
- Prevents the UI from tracking vendor upgrades

#### Alternatives considered

**Show raw vendor fields "for developers"**

- Optional later behind a debug flag, not MVP default

#### Consequences

- Copy review is part of `frontend-reviewer`
- Status mapping lives in Go, not in React `switch` statements that reimplement the adapter

---

## 3. Planned routes (MVP)

Align names with PRD menus. Exact paths TBD at implementation:

- Login (or SSO bounce)
- Dashboard (summary, trends, errors as the API provides)
- Use Cases list + detail

## 4. Out of scope

- Calling Orchestrator from the browser
- Robot start/stop controls
- Duplicating Go metric formulas in the client

## 5. Open questions

1. Host split (same origin vs API subdomain) and cookie `Domain`
2. Whether enterprise UI kit replaces shadcn
3. Polling interval vs scheduled snapshot (PRD freshness SLA)

## 6. References

- `docs/PRD.md` §17, §34, §35, §39
- `docs/adrs/ADR_BACKEND.md`
- `AGENTS.md`
