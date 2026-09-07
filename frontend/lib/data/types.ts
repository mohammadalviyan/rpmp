export type AutomationType = "Attended" | "Unattended";
export type UseCaseStatus = "Running" | "Warning" | "Stopped";
export type ErrorType = "Error A" | "Error B";
export type IssueStatus = "Open" | "In Progress" | "Resolved";
export type EmailStatus = "Sent" | "Pending" | "Failed";

export type Issue = {
  id: string;
  name: string;
  description: string;
  occurredAt: string;
  status: IssueStatus;
  errorType: ErrorType;
};

export type TrendPoint = {
  day: string;
  success: number;
  failed: number;
};

export type UseCase = {
  id: string;
  name: string;
  owner: string;
  description: string;
  successRate: number;
  status: UseCaseStatus;
  automationType: AutomationType;
  totalProcess: number;
  totalSuccess: number;
  totalFailed: number;
  issues: Issue[];
  trend: TrendPoint[];
};

export type Summary = {
  totalUseCase: number;
  successRate: number;
  totalIssue: number;
  unattended: number;
  attended: number;
  reportSent: number;
};

export type PerformancePoint = {
  month: string;
  success: number;
  failed: number;
  rate: number;
};

export type DistributionPoint = {
  name: string;
  value: number;
};

export type ReportRecord = {
  id: string;
  useCase: string;
  type: string;
  period: string;
  generatedAt: string;
  size: string;
  status: "Ready" | "Generating" | "Failed";
};

export type EmailRecord = {
  id: string;
  useCase: string;
  recipient: string;
  subject: string;
  date: string;
  status: EmailStatus;
};

export interface RpaDataProvider {
  getUseCases(): Promise<UseCase[]>;
  getUseCase(id: string): Promise<UseCase | undefined>;
  getSummary(): Promise<Summary>;
  getPerformanceSeries(): Promise<PerformancePoint[]>;
  getErrorDistribution(): Promise<DistributionPoint[]>;
  getAutomationTypeSplit(): Promise<DistributionPoint[]>;
  getReportHistory(): Promise<ReportRecord[]>;
  getEmailHistory(): Promise<EmailRecord[]>;
  getRecipients(): Promise<string[]>;
  getReportTypes(): Promise<string[]>;
  getReportPeriods(): Promise<string[]>;
}
