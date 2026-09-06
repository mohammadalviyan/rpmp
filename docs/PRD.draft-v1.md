# PRD — RPA Performance & Monitoring Platform

**Document Status:** Draft  
**Version:** 1.0  
**Date:** 2026-09-06  
**Product Name:** RPA Performance & Monitoring Platform  
**Document Type:** Product Requirements Document

---

## 1. Executive Summary

RPA Performance & Monitoring Platform adalah platform internal untuk mengubah data operasional RPA yang berasal dari Orchestrator menjadi informasi yang mudah dipahami oleh tim RPA, operational user, business user, dan management.

Platform ini dirancang untuk mengatasi dua permasalahan utama:

1. Proses reporting dan KRI masih membutuhkan pengambilan data manual dari database, penyusunan data, formula, dan template spreadsheet.
2. Monitoring RPA melalui UI Orchestrator cenderung berorientasi teknis sehingga kurang nyaman digunakan oleh business/user untuk memahami performa automation secara cepat.

MVP akan menyediakan dua menu utama:

- **Dashboard** — high-level overview kondisi dan performa RPA.
- **Use Cases** — daftar automation/use case beserta filter, status, health, dan detail performanya.

Produk tidak diposisikan sebagai pengganti Orchestrator, melainkan sebagai **information layer** di atas data RPA yang menerjemahkan raw operational data menjadi insight yang human-readable dan actionable.

> **Product Vision:** Turn raw RPA operational data into actionable performance and monitoring insights.

---

# 2. Background & Context

Data transaksi dan execution RPA pada dasarnya sudah tersedia di environment Orchestrator dan/atau database terkait. Namun, data tersebut belum secara langsung tersedia dalam format yang siap digunakan untuk reporting dan monitoring bisnis.

Current workflow secara umum:

```text
Orchestrator
     │
     ▼
Database / Transaction Data
     │
     │ Manual Query
     ▼
Spreadsheet
     │
     │ Manual Transformation / Formula
     ▼
KRI / Operational Report
     │
     ▼
User / Management
```

Di sisi lain, business/user yang ingin mengetahui kondisi automation harus berinteraksi langsung dengan Orchestrator yang memiliki istilah dan tampilan lebih teknis.

Platform ini akan menyediakan abstraction layer:

```text
Orchestrator / RPA Data
          │
          ▼
   RPA Data Layer
          │
          ▼
RPA Performance & Monitoring Platform
          │
     ┌────┴────┐
     ▼         ▼
 Dashboard   Use Cases
```

---

# 3. Problem Statement

## 3.1 Manual Reporting

Operational/RPA team masih perlu:

- Menjalankan query database secara manual.
- Mengambil data transaction/execution.
- Membersihkan dan menyusun data.
- Membuat formula pada spreadsheet.
- Memasukkan hasil ke template reporting.
- Mengulang proses tersebut pada periode berikutnya.

Akibatnya:

- Memerlukan waktu.
- Berpotensi terjadi human error.
- Logic reporting dapat berbeda antar pengguna.
- Historical comparison kurang praktis.
- Reporting sangat bergantung pada spreadsheet.

## 3.2 Technical Monitoring Experience

Orchestrator menyediakan informasi yang lengkap untuk kebutuhan teknis, tetapi belum tentu optimal untuk business-facing monitoring.

Business/user lebih membutuhkan jawaban seperti:

- Berapa automation yang aktif?
- Berapa success rate RPA saat ini?
- Use case mana yang bermasalah?
- Apakah failure meningkat?
- Error apa yang paling sering terjadi?
- Apakah automation memenuhi SLA?
- Bagaimana performanya dibanding periode sebelumnya?

## 3.3 Data Exists, But Insight Does Not

Masalah utama bukan ketiadaan data, melainkan:

> Data operasional tersedia, tetapi belum diubah menjadi information layer yang mudah dipahami dan digunakan untuk pengambilan keputusan.

---

# 4. Product Vision

## Vision

> **Create a single, human-readable source of truth for RPA performance, health, and operational monitoring.**

## Mission

Mengurangi pekerjaan manual dalam monitoring dan reporting RPA dengan menyediakan informasi performa automation secara terstruktur, visual, konsisten, dan actionable.

---

# 5. Product Positioning

Produk ini **bukan**:

- Pengganti Orchestrator.
- Tool untuk membuat atau menjalankan robot.
- Tool development/debugging RPA secara langsung.
- Database transaction viewer mentah.

Produk ini adalah:

> **RPA Performance & Monitoring Platform** yang berada di atas data operasional RPA untuk menyediakan monitoring, analytics, health indicators, dan reporting yang lebih mudah dipahami.

---

# 6. Goals & Objectives

## 6.1 MVP Goals

1. Menyediakan single view untuk performa RPA.
2. Menampilkan total dan status use case.
3. Menampilkan execution volume dan success/failure rate.
4. Mempermudah identifikasi use case yang membutuhkan attention.
5. Menyediakan detail performa per use case.
6. Mengurangi kebutuhan akses langsung ke Orchestrator untuk kebutuhan monitoring umum.
7. Menjadi foundation untuk automated KRI reporting di fase berikutnya.

## 6.2 Business Objectives

Target outcome:

- Mengurangi manual effort dalam reporting.
- Mengurangi ketergantungan terhadap spreadsheet.
- Mempercepat identifikasi automation yang bermasalah.
- Meningkatkan visibility terhadap portfolio RPA.
- Menyediakan data yang konsisten untuk operational monitoring dan reporting.

---

# 7. Non-Goals for MVP

Fitur berikut tidak menjadi scope MVP:

- Robot development.
- Robot deployment.
- Robot scheduling/control.
- Start/stop job dari dashboard.
- Direct transaction editing.
- Automated root-cause analysis berbasis AI.
- Predictive failure.
- Full KRI report automation.
- Advanced alerting/notification.
- Multi-tenant architecture.
- Self-service use case onboarding.
- Cost/ROI analytics yang kompleks.

Fitur tersebut dapat masuk roadmap fase berikutnya.

---

# 8. Target Users

## 8.1 RPA Developer

Kebutuhan:

- Melihat performance automation.
- Mengetahui failure/error yang sering terjadi.
- Melihat trend execution.
- Melakukan initial troubleshooting.

## 8.2 RPA Operator / Support

Kebutuhan:

- Monitoring kondisi automation.
- Identifikasi use case bermasalah.
- Melihat execution/failure trend.
- Menentukan prioritas follow-up.

## 8.3 Business User

Kebutuhan:

- Melihat kondisi use case secara sederhana.
- Memahami success rate.
- Mengetahui apakah automation sedang healthy.
- Tidak perlu memahami terminology teknis Orchestrator.

## 8.4 RPA Manager / Management

Kebutuhan:

- Portfolio overview.
- Performance summary.
- Health overview.
- Trend.
- High-level operational risk.

---

# 9. User Personas

| Persona | Primary Need | Technical Level |
|---|---|---|
| RPA Developer | Performance & error diagnostics | High |
| RPA Operator | Monitoring & incident detection | Medium |
| Business User | Human-readable status | Low–Medium |
| RPA Manager | Portfolio & performance overview | Medium |
| Management | Executive-level visibility | Low |

---

# 10. Product Scope

MVP terdiri dari:

```text
RPA Performance & Monitoring Platform
│
├── Authentication
│
├── Dashboard
│   ├── KPI Summary
│   ├── Execution Trend
│   ├── Success / Failure
│   ├── Use Case Status
│   ├── Top Errors
│   └── Attention Required
│
└── Use Cases
    ├── Search
    ├── Filter
    ├── Status
    ├── Health
    ├── Use Case List
    └── Use Case Detail
        ├── Summary
        ├── Performance
        ├── Execution Trend
        ├── Error Analysis
        └── Recent Activity
```

---

# 11. Information Architecture

```text
Login
  │
  ▼
Dashboard
  │
  ├───────────────┐
  ▼               ▼
Use Cases       Future
  │             Modules
  ▼
Use Case Detail
```

Future:

```text
Dashboard
Use Cases
KRI Reporting
Alerts
Analytics
Administration
```

---

# 12. Functional Requirements

# 12.1 Authentication

MVP membutuhkan authentication untuk internal user.

### MVP Requirement

User login menggunakan:

- Personal Number / Employee ID.
- Password.

### Minimum Roles

#### Viewer

- View Dashboard.
- View Use Cases.
- View Use Case Detail.

#### Admin

- Semua akses Viewer.
- Configuration management.
- User/role management apabila diperlukan.

### Future

Authentication dapat diintegrasikan dengan enterprise identity provider/SSO apabila tersedia.

---

# 12.2 Dashboard

Dashboard menjadi landing page utama.

## 12.2.1 Dashboard Header

Menampilkan:

- Product name.
- Dashboard title.
- Last updated timestamp.
- Selected period.
- Optional refresh action.

Contoh:

```text
RPA Performance & Monitoring

Overview of automation performance
Last updated: 06 Sep 2026 07:00
Period: Last 30 Days
```

---

## 12.2.2 KPI Cards

Minimum KPI:

### Total Use Cases

Menampilkan jumlah use case yang terdaftar.

### Active Use Cases

Jumlah use case aktif.

### Execution Volume

Jumlah execution/transaction pada selected period.

### Success Rate

Persentase successful execution.

Formula:

```text
Success Rate =
Successful Executions / Total Executions × 100%
```

### Failed Executions

Jumlah execution yang gagal.

### Optional KPI

- Average Execution Duration.
- SLA Achievement.
- Health Score.

---

# 12.3 Execution Trend

Dashboard menampilkan trend execution berdasarkan waktu.

Minimum:

- Daily aggregation.
- Success.
- Failed.
- Total.

Contoh:

```text
Execution Volume

│
│             ╭──╮
│         ╭───╯  ╰───╮
│    ╭────╯           ╰──
│────╯
└────────────────────────
  1   5   10  15  20  25  30
```

User dapat melihat:

- Apakah volume meningkat?
- Apakah terdapat abnormal drop?
- Apakah terdapat spike?

---

# 12.4 Success vs Failure

Menampilkan proporsi:

- Success.
- Failed.

Contoh:

```text
Success Rate
98.72%

Success       18,185
Failed           236
```

---

# 12.5 Use Case Status Overview

Minimum status:

- Active.
- Inactive.
- Error.
- Maintenance.

Status mapping perlu didefinisikan dalam data layer berdasarkan source Orchestrator.

---

# 12.6 Top Errors

Dashboard menampilkan error yang paling sering terjadi.

Contoh:

```text
Top Errors

SAP Timeout             42%
Invalid Data             27%
Login Failure            18%
Application Error        13%
```

Error perlu dinormalisasi apabila raw error dari Orchestrator terlalu granular.

Contoh:

```text
Raw:
Timeout while waiting for element...
Timeout while waiting for SAP...
SAP response timeout...

Normalized:
SAP Timeout
```

---

# 12.7 Attention Required

Dashboard harus menyediakan area untuk menunjukkan use case yang membutuhkan perhatian.

Contoh:

```text
3 Automations Require Attention

🔴 Report Generation
   Success rate dropped to 87.4%

🟡 Reconciliation
   Failure increased 28% this week

🟡 Invoice Processing
   Average duration increased 42%
```

### MVP Rule

Attention dapat ditentukan menggunakan rule-based logic.

Contoh:

```text
Critical:
Success Rate < 90%

Warning:
Success Rate >= 90% AND < 97%

Healthy:
Success Rate >= 97%
```

Threshold harus configurable di future.

---

# 12.8 Use Case List

Menu Use Cases menampilkan seluruh automation/use case.

## Filter

Minimum:

- Search.
- Business Unit / Department.
- Status.
- Health.
- Date / Period.

## Table Columns

Minimum:

| Field | Description |
|---|---|
| Use Case | Human-readable use case name |
| Business Unit | Owner/business unit |
| Status | Current status |
| Execution | Execution volume |
| Success Rate | Success percentage |
| Last Run | Last execution |
| Health | Current health indicator |

---

# 12.9 Use Case Detail

Ketika user memilih use case:

```text
← Back to Use Cases

LC Amendment
Trade Finance

● Active
Health Score: 92 / 100
```

## Summary Metrics

Minimum:

- Total executions.
- Successful executions.
- Failed executions.
- Success rate.
- Average duration.
- Last execution.
- SLA achievement, jika tersedia.

---

# 12.10 Use Case Performance Trend

Menampilkan:

- Success rate trend.
- Execution volume.
- Failure volume.
- Average duration.

User dapat memilih period:

- Today.
- Last 7 Days.
- Last 30 Days.
- Custom Period — future/optional.

---

# 12.11 Use Case Error Analysis

Menampilkan:

- Most common errors.
- Error count.
- Error percentage.
- Trend.
- Last occurrence.

Contoh:

```text
Most Common Errors

SAP Timeout       18   42%
Invalid Data      11   26%
Login Failure      7   16%
Others             7   16%
```

---

# 12.12 Recent Activity

Menampilkan execution terbaru:

| Time | Status | Duration | Error |
|---|---|---|---|
| 07:02 | Success | 03:41 | - |
| 06:58 | Failed | 02:13 | SAP Timeout |
| 06:42 | Success | 03:50 | - |

Detail transaction-level dapat menjadi future scope apabila diperlukan.

---

# 13. Health Model

Health merupakan abstraction layer untuk membuat kondisi automation lebih mudah dipahami.

Minimum:

```text
🟢 Healthy
🟡 Attention
🔴 Critical
```

Contoh rule MVP:

| Health | Rule |
|---|---|
| Healthy | Success Rate >= 97% |
| Attention | 90% <= Success Rate < 97% |
| Critical | Success Rate < 90% |

Rule ini merupakan starting point dan harus divalidasi dengan actual operational baseline.

---

# 14. Health Score — Future Enhancement

Pada fase berikutnya, health tidak hanya berdasarkan success rate.

Contoh:

```text
Health Score =
    Success Rate
    + Failure Trend
    + SLA Achievement
    + Availability
    + Execution Duration
```

Contoh:

```text
LC Amendment

Health Score
██████████████████░░ 92

Success Rate        98.7%
Failure Rate         1.3%
SLA                  98.9%
Availability        99.9%

Status: Healthy
```

---

# 15. Business Metadata

Platform harus memiliki kemampuan untuk menyimpan metadata yang lebih business-friendly dibanding raw Orchestrator data.

Contoh:

| Metadata | Example |
|---|---|
| Technical Process Name | CR_IMPORT_LC_AMEND |
| Display Name | LC Amendment |
| Business Process | Trade Finance |
| Business Unit | Trade Operations |
| Business Owner | Trade Operations |
| RPA Owner | RPA Team |
| Criticality | High |
| SLA | 08:00–17:00 |
| Expected Volume | 2,000/day |

Business metadata dapat menjadi master/configuration layer terpisah dari data Orchestrator.

---

# 16. Data Architecture

## 16.1 Recommended Architecture

```text
                 ┌────────────────────┐
                 │    Orchestrator    │
                 └─────────┬──────────┘
                           │
                 API / DB / Other Source
                           │
                           ▼
                 ┌────────────────────┐
                 │   Data Collector   │
                 └─────────┬──────────┘
                           │
                           ▼
                 ┌────────────────────┐
                 │   RPA Data Layer   │
                 └─────────┬──────────┘
                           │
                           ▼
                 ┌────────────────────┐
                 │ Backend / API      │
                 └─────────┬──────────┘
                           │
                           ▼
                 ┌────────────────────┐
                 │ Frontend Dashboard │
                 └────────────────────┘
```

## 16.2 Data Source Strategy

Prioritas eksplorasi:

1. Official Orchestrator API.
2. Approved database access/read replica.
3. Other supported data export/integration.
4. Direct production DB access only if formally approved.

Dashboard tidak boleh bergantung langsung pada raw production database.

---

# 17. Data Abstraction Layer

Data layer harus memisahkan:

```text
Source Data
    ↓
Raw Model
    ↓
Normalized Model
    ↓
Business Metrics
    ↓
Dashboard API
```

Contoh:

```text
Orchestrator:
JobState = 3

Data Layer:
ExecutionStatus = FAILED

Dashboard:
Failed Execution
```

Hal ini menjaga frontend agar tidak tightly coupled dengan vendor-specific terminology.

---

# 18. Proposed Backend API

Contoh endpoint MVP:

```http
GET /api/v1/dashboard/summary

GET /api/v1/dashboard/execution-trend

GET /api/v1/dashboard/errors

GET /api/v1/dashboard/attention

GET /api/v1/use-cases

GET /api/v1/use-cases/{id}

GET /api/v1/use-cases/{id}/summary

GET /api/v1/use-cases/{id}/performance

GET /api/v1/use-cases/{id}/errors

GET /api/v1/use-cases/{id}/executions
```

API contract perlu didefinisikan setelah discovery terhadap actual Orchestrator data.

---

# 19. Data Model — Initial Concept

## Use Case

```text
UseCase
---------
id
technical_name
display_name
description
business_unit
business_owner
rpa_owner
status
criticality
sla
created_at
updated_at
```

## Execution

```text
Execution
---------
id
use_case_id
execution_id
started_at
finished_at
duration
status
error_code
error_message
transaction_count
created_at
```

## Error

```text
Error
---------
id
use_case_id
category
error_code
message
occurred_at
count
```

## Business Metadata

```text
BusinessMetadata
----------------
use_case_id
business_process
business_unit
business_owner
criticality
sla
expected_volume
```

---

# 20. Data Freshness

Dashboard harus menampilkan timestamp data terakhir.

Contoh:

```text
Last updated: 06 Sep 2026 07:00
```

MVP dapat menggunakan batch refresh.

Target awal:

- Refresh interval: configurable.
- Default: 5–15 minutes, tergantung capability source.
- Historical data dapat diproses asynchronously.

Real-time monitoring bukan mandatory untuk MVP kecuali dibutuhkan operationally.

---

# 21. Data Accuracy

Data pada platform harus dapat direkonsiliasi dengan source Orchestrator.

Target:

> **>= 99.5% data accuracy**

Minimal reconciliation:

```text
Orchestrator execution count
          vs
Platform execution count
```

Periode rekonsiliasi perlu ditentukan dalam technical design.

---

# 22. Security Requirements

MVP bersifat internal-only.

Minimum requirements:

- Authentication.
- Password tidak disimpan dalam plaintext.
- HTTPS.
- Authorization berdasarkan role.
- Audit logging untuk administrative actions.
- No direct frontend access to Orchestrator DB.
- Sensitive credentials disimpan menggunakan approved secret management.
- Least privilege database/API access.

---

# 23. RBAC

Recommended initial model:

```text
Role
├── Viewer
│   ├── Dashboard
│   ├── Use Cases
│   └── Use Case Detail
│
└── Admin
    ├── All Viewer permissions
    └── Configuration
```

Future:

```text
Super Admin
RPA Admin
RPA Developer
RPA Operator
Business User
Management
```

Future authorization dapat berdasarkan Business Unit:

```text
Business User
     ↓
Only see assigned business unit
```

---

# 24. UX Principles

## 24.1 Human-readable

Hindari terminology teknis Orchestrator apabila tidak diperlukan.

Gunakan:

> Success Rate

daripada:

> JobState = Successful

## 24.2 Action-oriented

Dashboard harus menjawab:

> What happened?

> Why?

> What needs attention?

## 24.3 Progressive Disclosure

Level informasi:

```text
Level 1 — Portfolio
   ↓
Dashboard

Level 2 — Problematic Automation
   ↓
Use Case List

Level 3 — Specific Automation
   ↓
Use Case Detail

Level 4 — Execution
   ↓
Recent Activity / Transaction
```

## 24.4 Avoid Chart Overload

Dashboard bukan kumpulan chart.

Prioritas:

1. KPI.
2. Health.
3. Trend.
4. Attention.
5. Detail.

---

# 25. Dashboard UX Concept

```text
┌─────────────────────────────────────────────────────────┐
│ RPA Performance & Monitoring                            │
│ Overview of automation performance                      │
│ Last updated: 07:00                 Period: Last 30 Days│
├─────────────────────────────────────────────────────────┤
│                                                         │
│ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ │
│ │Use Case│ │ Active │ │Execute │ │Success │ │ Failed │ │
│ │   42   │ │   35   │ │18,421  │ │ 98.72% │ │  127   │ │
│ └────────┘ └────────┘ └────────┘ └────────┘ └────────┘ │
│                                                         │
│ ┌──────────────────────────┐ ┌────────────────────────┐ │
│ │ Execution Trend          │ │ Success / Failure      │ │
│ │                          │ │                        │ │
│ │          GRAPH           │ │       GRAPH            │ │
│ └──────────────────────────┘ └────────────────────────┘ │
│                                                         │
│ ┌──────────────────────────┐ ┌────────────────────────┐ │
│ │ Top Errors               │ │ Attention Required     │ │
│ │                          │ │                        │ │
│ │ SAP Timeout       42%    │ │ 🔴 Report Generation   │ │
│ │ Invalid Data       27%   │ │ 🟡 Reconciliation     │ │
│ │ Login Failure      18%   │ │ 🟡 Invoice Processing  │ │
│ └──────────────────────────┘ └────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

---

# 26. Use Case UX Concept

```text
┌─────────────────────────────────────────────────────────┐
│ Use Cases                                               │
│                                                         │
│ Search [________________]                               │
│                                                         │
│ Business Unit [All ▼]  Status [All ▼] Health [All ▼]  │
├─────────────────────────────────────────────────────────┤
│ Use Case           Status   Execution Success Health   │
│                                                         │
│ Invoice Processing  Active   4,218     99.2%    🟢      │
│ LC Amendment        Active   2,341     98.7%    🟢      │
│ Reconciliation      Active   1,842     95.1%    🟡      │
│ Report Generation   Error      421     87.4%    🔴      │
└─────────────────────────────────────────────────────────┘
```

---

# 27. Use Case Detail UX Concept

```text
┌─────────────────────────────────────────────────────────┐
│ ← Back to Use Cases                                     │
│                                                         │
│ LC Amendment                                            │
│ Trade Finance                                           │
│ ● Active                          Health: 92 / 100      │
├─────────────────────────────────────────────────────────┤
│                                                         │
│ Executions │ Success Rate │ Avg Duration │ Failed       │
│   2,341    │    98.7%     │    03:42     │   31         │
│                                                         │
├──────────────────────────────┬──────────────────────────┤
│ Performance Trend            │ Error Analysis           │
│                              │                          │
│           GRAPH              │ SAP Timeout        42%   │
│                              │ Invalid Data        27%  │
│                              │ Login Failure       18%  │
├──────────────────────────────┴──────────────────────────┤
│ Recent Activity                                         │
│                                                         │
│ Time    Status       Duration      Error                │
│ 07:02   Success      03:41         -                    │
│ 06:58   Failed       02:13         SAP Timeout          │
└─────────────────────────────────────────────────────────┘
```

---

# 28. Reporting & KRI Roadmap

KRI automation bukan MVP core feature, tetapi merupakan salah satu target utama product evolution.

## Current

```text
Database
 ↓
Manual Query
 ↓
Excel
 ↓
Formula
 ↓
KRI Report
```

## Future

```text
RPA Data Layer
 ↓
Metric Engine
 ↓
KRI Calculation
 ↓
KRI Dashboard
 ↓
Excel / PDF
 ↓
Scheduled Report
```

Potential future metrics:

- Total executions.
- Failed executions.
- Success rate.
- Critical failures.
- SLA achievement.
- Availability.
- Average processing duration.
- Volume processed.
- Exception rate.

---

# 29. Future Roadmap

## Phase 0 — Discovery

Objectives:

- Understand Orchestrator capabilities.
- Explore API.
- Understand DB schema.
- Identify available execution data.
- Identify error/log data.
- Define data dictionary.
- Validate security/access model.

Deliverables:

```text
Orchestrator Data Inventory
Data Dictionary
Integration Assessment
Initial Architecture
MVP Scope Validation
```

---

## Phase 1 — MVP

### Dashboard

- KPI cards.
- Execution trend.
- Success/failure.
- Use case status.
- Top errors.
- Attention required.

### Use Cases

- List.
- Search.
- Filter.
- Status.
- Health.
- Detail.
- Performance trend.
- Error analysis.
- Recent activity.

### Platform

- Authentication.
- Basic RBAC.
- Backend API.
- Data abstraction layer.

---

## Phase 2 — Operational Intelligence

- Health Score.
- SLA monitoring.
- Duration anomaly.
- Failure trend.
- Better error categorization.
- Advanced attention rules.
- Availability monitoring.
- Operational alerts.

---

## Phase 3 — Reporting Automation

- KRI dashboard.
- Automated KRI calculation.
- Excel export.
- PDF export.
- Scheduled reports.
- Reporting templates.
- Historical comparison.

---

## Phase 4 — Business Intelligence

- Automation ROI.
- Hours saved.
- Transaction volume.
- FTE equivalent.
- Cost analytics.
- Business impact.
- Portfolio analytics.

---

# 30. Success Metrics

MVP success should be measured by outcomes, not only feature completion.

## Operational Efficiency

Target:

> **Reduce manual reporting effort by >= 70%.**

## Monitoring Speed

Target:

> User can identify problematic use cases within **< 1 minute**.

## Data Accuracy

Target:

> **>= 99.5%** reconciliation accuracy against source data.

## Adoption

Target:

> **>= 80%** of target operational users use the platform for routine monitoring.

## Reporting Dependency

Target:

> Reduce routine dependency on manual spreadsheet transformation.

## Availability

Target:

> **>= 99%** platform availability during agreed operational hours.

---

# 31. MVP Acceptance Criteria

## Dashboard

- [ ] User can access dashboard after login.
- [ ] Dashboard shows total use cases.
- [ ] Dashboard shows active use cases.
- [ ] Dashboard shows execution volume.
- [ ] Dashboard shows success rate.
- [ ] Dashboard shows failed executions.
- [ ] Dashboard shows execution trend.
- [ ] Dashboard shows success/failure distribution.
- [ ] Dashboard shows top errors.
- [ ] Dashboard shows attention-required use cases.
- [ ] Dashboard displays data freshness.

## Use Cases

- [ ] User can open Use Cases menu.
- [ ] User can search use cases.
- [ ] User can filter by status.
- [ ] User can filter by business unit.
- [ ] User can filter by health.
- [ ] User can view execution count.
- [ ] User can view success rate.
- [ ] User can view last execution.
- [ ] User can open use case detail.

## Use Case Detail

- [ ] User can view summary metrics.
- [ ] User can view success rate.
- [ ] User can view execution trend.
- [ ] User can view error analysis.
- [ ] User can view recent activity.

## Security

- [ ] Authentication is required.
- [ ] Viewer role is implemented.
- [ ] Admin role is implemented where required.
- [ ] Passwords are securely stored.
- [ ] Frontend does not directly access production DB.

---

# 32. Risks & Mitigation

| Risk | Impact | Mitigation |
|---|---|---|
| Orchestrator API unavailable | High | Explore API early; prepare approved DB integration |
| Production DB cannot be exposed | High | Introduce data collector/read replica/service layer |
| DB schema changes | Medium | Abstract source data behind data layer |
| Data mismatch | High | Reconciliation process |
| Error messages too granular | Medium | Error normalization |
| Business metadata unavailable | Medium | Maintain separate metadata configuration |
| Authentication integration delayed | Medium | Start with internal auth, prepare SSO interface |
| Dashboard becomes too complex | Medium | Keep MVP focused |
| Real-time data expectation | Medium | Define explicit data freshness SLA |

---

# 33. Key Technical Decisions to Validate

These questions should be answered during Phase 0:

### Orchestrator

- What Orchestrator product/version is used?
- Does it expose official APIs?
- Which API endpoints provide execution history?
- Which API endpoints provide job status?
- Can error information be retrieved?
- What are API rate limits?
- What authentication mechanism is available?

### Database

- What DB engine is used?
- Is direct read access allowed?
- Is a read replica available?
- Which tables contain execution history?
- Which tables contain error information?
- How much historical data exists?
- What is the retention policy?

### Data Volume

- Daily execution volume?
- Monthly execution volume?
- Number of use cases?
- Average transaction size?
- Expected growth?

### Security

- Can service account be created?
- Is SSO available?
- Which enterprise identity provider is used?
- What audit requirements apply?

---

# 34. Recommended MVP Technology Direction

Technology should follow the organization's existing standard, but an initial architecture could be:

```text
Frontend
React / Next.js

Backend
Go

Database
PostgreSQL / approved relational database

Integration
Orchestrator API
        OR
Approved read-only DB access

Authentication
Internal authentication initially
SSO-ready architecture

Deployment
Internal infrastructure / approved platform
```

The technology stack is not locked by this PRD and must be validated against enterprise standards.

---

# 35. Non-Functional Requirements

## Performance

Target:

- Dashboard initial load: < 3 seconds under normal conditions.
- API response for common queries: < 1 second where cached/indexed.
- Large historical analytics should use aggregation/precomputation where necessary.

## Availability

Target:

- >= 99% during agreed operational hours.

## Scalability

Architecture should support growth in:

- Number of use cases.
- Execution volume.
- Historical data.
- Number of users.

## Maintainability

- Clear separation between source integration, business logic, and presentation.
- API-first architecture.
- Versioned APIs.
- Centralized logging.
- Automated tests.

---

# 36. Observability

Platform should eventually monitor itself.

Minimum:

- Application logs.
- API logs.
- Integration failures.
- Data synchronization status.
- Data freshness.
- Authentication failures.

Future:

```text
Platform Health

Data Pipeline       🟢
Orchestrator API    🟢
Database            🟢
Backend             🟢
Frontend             🟢
```

---

# 37. Auditability

Important actions should be logged:

- Login.
- Logout.
- Configuration changes.
- User/role changes.
- Metadata changes.
- Administrative actions.

Future:

```text
Audit Log

User       Action              Time
12345      Updated Use Case    07:31
67890      Changed SLA         08:02
```

---

# 38. Data Governance

Data ownership should be clearly defined.

Suggested ownership:

```text
Orchestrator Data
      │
      ▼
RPA / Platform Team
      │
      ▼
Business Metadata
      │
      ▼
Business / RPA Owner
```

Each use case should ideally have:

- Business owner.
- RPA owner.
- Criticality.
- SLA.
- Business unit.

This makes future portfolio analytics possible.

---

# 39. Product Principles

The platform should follow these principles:

### 1. Human-readable first

Technical data must be translated into business-friendly information.

### 2. One source of truth

Metrics should be calculated consistently.

### 3. Action over information

The platform should highlight what needs attention.

### 4. Separate source from presentation

Orchestrator-specific implementation details should not leak into the frontend.

### 5. MVP first

Do not build advanced analytics before validating core monitoring behavior.

### 6. Automation of manual work

The product should continuously reduce spreadsheet/manual processes.

---

# 40. Future Feature Ideas

Potential features after MVP:

## Smart Alerts

```text
LC Amendment success rate dropped below 90%.
```

## Anomaly Detection

```text
Execution duration is 45% higher than normal.
```

## SLA Monitoring

```text
12 use cases breached SLA today.
```

## Incident Management

```text
Create incident
Assign owner
Track resolution
```

## Root Cause Analytics

```text
Most failures originated from:
SAP → 62%
Network → 21%
Data → 17%
```

## Business Impact

```text
Transactions processed     1.2M
Hours saved                 18,421
Estimated FTE equivalent   11.4
```

## Portfolio View

```text
42 Use Cases
35 Healthy
5 Attention
2 Critical
```

## Scheduled Reporting

```text
Every month:
Generate KRI report
→ Excel
→ PDF
→ Distribution list
```

---

# 41. Product Evolution

The intended evolution is:

```text
                   RPA PERFORMANCE &
                  MONITORING PLATFORM
                           │
             ┌─────────────┼─────────────┐
             ▼             ▼             ▼
        Monitoring      Reporting      Analytics
             │             │             │
             ▼             ▼             ▼
         Dashboard        KRI          Business
         Use Cases       Reports       Impact
             │             │             │
             └─────────────┼─────────────┘
                           ▼
                    RPA Command Center
```

Long-term vision:

> A centralized operational intelligence platform for the organization's entire RPA portfolio.

---

# 42. Final MVP Definition

The MVP is considered successful when an internal user can:

```text
Login
  ↓
Open Dashboard
  ↓
Understand overall RPA health
  ↓
Identify problematic automation
  ↓
Open Use Cases
  ↓
Filter/Search a specific automation
  ↓
Open Use Case Detail
  ↓
Understand performance
  ↓
Identify common errors
  ↓
Determine whether follow-up is required
```

without needing to access Orchestrator for routine monitoring.

---

# 43. One-Sentence Product Definition

> **RPA Performance & Monitoring Platform is an internal operational intelligence platform that transforms RPA execution data into human-readable performance, health, and monitoring insights while reducing manual reporting effort.**

---

# 44. Recommended Next Steps

Before development starts:

1. **Orchestrator Discovery**
   - Explore official API.
   - Map available endpoints.
   - Identify execution/job/error data.

2. **Data Discovery**
   - Map existing database tables.
   - Define source-of-truth.
   - Define data dictionary.

3. **Metric Definition**
   - Define Success Rate.
   - Define Failure.
   - Define Execution.
   - Define Active/Inactive.
   - Define Health.
   - Define SLA.

4. **UX Validation**
   - Create low-fidelity wireframes.
   - Validate with RPA users.
   - Validate with business users.

5. **Technical Architecture**
   - Finalize integration approach.
   - Define backend API.
   - Define normalized data model.

6. **MVP Development**
   - Authentication.
   - Data ingestion.
   - Backend.
   - Dashboard.
   - Use Cases.
   - Use Case Detail.

7. **Pilot**
   - Select several representative use cases.
   - Compare dashboard metrics against existing manual reports.
   - Measure time saved.
   - Collect user feedback.

---

# Appendix A — Example KPI Definitions

| KPI | Definition |
|---|---|
| Total Use Cases | Number of registered use cases |
| Active Use Cases | Use cases currently marked active |
| Execution Volume | Total execution/transaction count in selected period |
| Successful Execution | Execution completed successfully |
| Failed Execution | Execution completed with failure/exception |
| Success Rate | Successful executions / total executions |
| Failure Rate | Failed executions / total executions |
| Average Duration | Average execution duration |
| SLA Achievement | Executions meeting defined SLA / applicable executions |
| Health | Rule-based abstraction of current automation condition |

---

# Appendix B — Example Business Questions

Dashboard should answer:

### Portfolio

- How many automations do we have?
- How many are active?
- What percentage are healthy?

### Performance

- How many executions happened?
- What is the success rate?
- Is performance improving or declining?

### Operations

- Which use cases need attention?
- What errors occur most frequently?
- Which use case has the highest failure rate?

### Business

- Which business unit has the most automation?
- Which automation processes the highest volume?
- Which automation is most critical?

### Future Management

- How much operational work is automated?
- How many hours are saved?
- What is the business value of the RPA portfolio?

---

# Appendix C — MVP vs Future

| Capability | MVP | Future |
|---|:---:|:---:|
| Authentication | ✓ | SSO |
| RBAC | Basic | Advanced |
| Dashboard | ✓ | Advanced |
| Use Case List | ✓ | Advanced |
| Use Case Detail | ✓ | Advanced |
| Success Rate | ✓ | ✓ |
| Error Analysis | ✓ | AI-assisted |
| Health | Basic | Advanced Score |
| SLA | Optional | ✓ |
| Attention Required | Basic | Smart |
| KRI | — | ✓ |
| Excel Export | — | ✓ |
| PDF Export | — | ✓ |
| Scheduled Report | — | ✓ |
| Alerts | — | ✓ |
| Anomaly Detection | — | ✓ |
| Root Cause Analysis | — | ✓ |
| ROI / Business Impact | — | ✓ |
| Command Center | — | Long-term |

---

**End of PRD**
