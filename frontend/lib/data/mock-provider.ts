import type {
  EmailRecord,
  ReportRecord,
  RpaDataProvider,
  Summary,
  UseCase,
} from "@/lib/data/types";

const trend = (base: number) =>
  ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"].map((day, index) => ({
    day,
    success: Math.round(base + Math.sin(index) * 18 + index * 4),
    failed: Math.max(1, Math.round(8 + Math.cos(index * 1.4) * 5)),
  }));

const useCases: UseCase[] = [
  {
    id: "use-case-a",
    name: "Use Case A",
    owner: "Finance Operations",
    description:
      "Automated invoice extraction, validation and posting to the ERP ledger.",
    successRate: 94,
    status: "Running",
    automationType: "Unattended",
    totalProcess: 1240,
    totalSuccess: 1166,
    totalFailed: 74,
    trend: trend(150),
    issues: [
      {
        id: "a-1",
        name: "Invoice parser timeout",
        description:
          "OCR engine did not respond within 30s while reading invoice batch #4471.",
        occurredAt: "03 Sep 2026, 09:14",
        status: "Open",
        errorType: "Error A",
      },
      {
        id: "a-2",
        name: "Duplicate ledger entry",
        description:
          "Record already posted for vendor PT Nusantara, transaction skipped.",
        occurredAt: "02 Sep 2026, 16:02",
        status: "In Progress",
        errorType: "Error B",
      },
      {
        id: "a-3",
        name: "Session lost",
        description:
          "ERP session expired mid-run, bot restarted the queue automatically.",
        occurredAt: "01 Sep 2026, 11:38",
        status: "Resolved",
        errorType: "Error A",
      },
    ],
  },
  {
    id: "use-case-b",
    name: "Use Case B",
    owner: "Customer Service",
    description:
      "Attended bot assisting agents with ticket enrichment and SLA tagging.",
    successRate: 88,
    status: "Running",
    automationType: "Attended",
    totalProcess: 860,
    totalSuccess: 757,
    totalFailed: 103,
    trend: trend(110),
    issues: [
      {
        id: "b-1",
        name: "Ticket field mismatch",
        description:
          "Priority field returned an unmapped value from the CRM API.",
        occurredAt: "03 Sep 2026, 08:20",
        status: "Open",
        errorType: "Error B",
      },
      {
        id: "b-2",
        name: "Credential rejected",
        description:
          "Agent workstation credential rotation broke the login step.",
        occurredAt: "02 Sep 2026, 13:45",
        status: "Resolved",
        errorType: "Error A",
      },
    ],
  },
  {
    id: "use-case-c",
    name: "Use Case C",
    owner: "Human Resources",
    description:
      "Daily attendance reconciliation and payroll pre-check reporting.",
    successRate: 76,
    status: "Warning",
    automationType: "Unattended",
    totalProcess: 640,
    totalSuccess: 486,
    totalFailed: 154,
    trend: trend(80),
    issues: [
      {
        id: "c-1",
        name: "Attendance file missing",
        description:
          "Source export was not published to the shared folder before the run.",
        occurredAt: "03 Sep 2026, 06:05",
        status: "Open",
        errorType: "Error A",
      },
      {
        id: "c-2",
        name: "Cleansing rule failed",
        description:
          "Unexpected date format dd-mm-yy caused 42 rows to be rejected.",
        occurredAt: "02 Sep 2026, 07:12",
        status: "In Progress",
        errorType: "Error B",
      },
      {
        id: "c-3",
        name: "Report template broken",
        description:
          "Pivot range shifted after a manual column insert in the template.",
        occurredAt: "31 Aug 2026, 18:24",
        status: "Open",
        errorType: "Error B",
      },
    ],
  },
  {
    id: "use-case-d",
    name: "Use Case D",
    owner: "Procurement",
    description:
      "Vendor master data cleansing and purchase order status monitoring.",
    successRate: 91,
    status: "Running",
    automationType: "Attended",
    totalProcess: 520,
    totalSuccess: 473,
    totalFailed: 47,
    trend: trend(70),
    issues: [
      {
        id: "d-1",
        name: "Vendor code not found",
        description:
          "Lookup failed for 6 vendors that were archived last quarter.",
        occurredAt: "02 Sep 2026, 10:41",
        status: "Resolved",
        errorType: "Error B",
      },
    ],
  },
  {
    id: "use-case-e",
    name: "Use Case E",
    owner: "Risk & Compliance",
    description:
      "Regulatory report compilation with automated distribution to reviewers.",
    successRate: 69,
    status: "Stopped",
    automationType: "Unattended",
    totalProcess: 310,
    totalSuccess: 214,
    totalFailed: 96,
    trend: trend(45),
    issues: [
      {
        id: "e-1",
        name: "Portal unavailable",
        description:
          "Regulator portal returned HTTP 503 during the submission step.",
        occurredAt: "03 Sep 2026, 05:30",
        status: "Open",
        errorType: "Error A",
      },
      {
        id: "e-2",
        name: "Checksum mismatch",
        description:
          "Generated report checksum differs from the validated baseline.",
        occurredAt: "01 Sep 2026, 21:15",
        status: "In Progress",
        errorType: "Error B",
      },
      {
        id: "e-3",
        name: "Queue overflow",
        description:
          "Pending items exceeded the configured queue limit of 500.",
        occurredAt: "30 Aug 2026, 14:50",
        status: "Open",
        errorType: "Error A",
      },
    ],
  },
];

const summary: Summary = {
  totalUseCase: 25,
  successRate: 85,
  totalIssue: 13,
  unattended: 14,
  attended: 11,
  reportSent: 20,
};

const performanceSeries = [
  { month: "Mar", success: 820, failed: 96, rate: 90 },
  { month: "Apr", success: 910, failed: 130, rate: 88 },
  { month: "May", success: 980, failed: 165, rate: 86 },
  { month: "Jun", success: 1120, failed: 140, rate: 89 },
  { month: "Jul", success: 1210, failed: 190, rate: 86 },
  { month: "Aug", success: 1340, failed: 210, rate: 86 },
  { month: "Sep", success: 1420, failed: 250, rate: 85 },
];

const errorDistribution = [
  { name: "Error A", value: 7 },
  { name: "Error B", value: 6 },
];

const automationTypeSplit = [
  { name: "Unattended", value: 14 },
  { name: "Attended", value: 11 },
];

const reportHistory: ReportRecord[] = [
  {
    id: "RPT-2091",
    useCase: "Use Case A",
    type: "Performance Summary",
    period: "Aug 2026",
    generatedAt: "04 Sep 2026, 06:00",
    size: "1.2 MB",
    status: "Ready",
  },
  {
    id: "RPT-2090",
    useCase: "Use Case C",
    type: "Issue Analysis",
    period: "Aug 2026",
    generatedAt: "03 Sep 2026, 06:00",
    size: "864 KB",
    status: "Ready",
  },
  {
    id: "RPT-2089",
    useCase: "Use Case B",
    type: "Daily Operations",
    period: "02 Sep 2026",
    generatedAt: "03 Sep 2026, 05:45",
    size: "402 KB",
    status: "Ready",
  },
  {
    id: "RPT-2088",
    useCase: "Use Case E",
    type: "Compliance Report",
    period: "Aug 2026",
    generatedAt: "03 Sep 2026, 05:30",
    size: "—",
    status: "Failed",
  },
  {
    id: "RPT-2087",
    useCase: "Use Case D",
    type: "Performance Summary",
    period: "Aug 2026",
    generatedAt: "02 Sep 2026, 06:00",
    size: "980 KB",
    status: "Ready",
  },
];

const emailHistory: EmailRecord[] = [
  {
    id: "EM-5512",
    useCase: "Use Case A",
    recipient: "finance.ops@company.com",
    subject: "Monthly Automation Report — Use Case A",
    date: "04 Sep 2026",
    status: "Sent",
  },
  {
    id: "EM-5511",
    useCase: "Use Case B",
    recipient: "cs.lead@company.com",
    subject: "Weekly Performance Digest — Use Case B",
    date: "03 Sep 2026",
    status: "Sent",
  },
  {
    id: "EM-5510",
    useCase: "Use Case C",
    recipient: "hr.report@company.com",
    subject: "Attendance Reconciliation Report",
    date: "03 Sep 2026",
    status: "Pending",
  },
  {
    id: "EM-5509",
    useCase: "Use Case E",
    recipient: "compliance@company.com",
    subject: "Regulatory Submission Summary",
    date: "02 Sep 2026",
    status: "Failed",
  },
  {
    id: "EM-5508",
    useCase: "Use Case D",
    recipient: "procurement@company.com",
    subject: "Vendor Data Quality Report",
    date: "02 Sep 2026",
    status: "Sent",
  },
];

const recipients = [
  "finance.ops@company.com",
  "cs.lead@company.com",
  "hr.report@company.com",
  "compliance@company.com",
  "procurement@company.com",
  "management@company.com",
];

const reportTypes = [
  "Performance Summary",
  "Issue Analysis",
  "Daily Operations",
  "Compliance Report",
];

const reportPeriods = [
  "Today",
  "Last 7 days",
  "Last 30 days",
  "Aug 2026",
  "Q3 2026",
];

export const mockDataProvider: RpaDataProvider = {
  async getUseCases() {
    return useCases;
  },
  async getUseCase(id) {
    return useCases.find((useCase) => useCase.id === id);
  },
  async getSummary() {
    return summary;
  },
  async getPerformanceSeries() {
    return performanceSeries;
  },
  async getErrorDistribution() {
    return errorDistribution;
  },
  async getAutomationTypeSplit() {
    return automationTypeSplit;
  },
  async getReportHistory() {
    return reportHistory;
  },
  async getEmailHistory() {
    return emailHistory;
  },
  async getRecipients() {
    return recipients;
  },
  async getReportTypes() {
    return reportTypes;
  },
  async getReportPeriods() {
    return reportPeriods;
  },
};
