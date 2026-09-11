export type UserRole = "viewer" | "admin";

export type User = {
  id: string;
  employee_id: string;
  display_name: string;
  role: UserRole;
};

export type AuthResponse = {
  user: User;
};

export type LoginRequest = {
  employee_id: string;
  password: string;
};

export type ApiError = {
  code: string;
  message: string;
};

export type SyncRunStatus = "running" | "success" | "failure" | "never";

export type SyncStatus = {
  run_id: string | null;
  status: SyncRunStatus;
  started_at: string | null;
  finished_at: string | null;
  rows_written: number;
};

export type SyncStarted = {
  run_id: string;
  status: "running";
};

export type DashboardSummary = {
  period: DashboardPeriod;
  freshness: {
    status: string;
    last_successful_refresh_at: string;
  };
  kpis: {
    total_use_cases: number;
    active_use_cases: number;
    execution_volume: number;
    success_rate: number | null;
    failed_executions: number;
  };
};

export type DashboardPeriod = {
  from: string;
  to: string;
  timezone: string;
};

export type ExecutionTrend = {
  period: DashboardPeriod;
  points: Array<{
    bucket: string;
    label: string;
    success: number;
    failure: number;
  }>;
};

export type DashboardErrors = {
  period: DashboardPeriod;
  groups: Array<{
    code: string;
    label: string;
    count: number;
  }>;
};

export type DataFreshness = {
  status: string;
  last_successful_refresh_at: string;
};

export type ApiUseCaseStatus = "active" | "inactive";

export type UseCaseListItem = {
  id: string;
  name: string;
  status: ApiUseCaseStatus;
  process_count: number;
  successful_count: number;
  error_count: number;
  stopped_count: number;
  failed_count: number;
  execution_volume: number;
  success_rate: number | null;
  environments: string[];
};

export type UseCasesResponse = {
  freshness: DataFreshness;
  items: UseCaseListItem[];
};

export type UseCaseProcess = {
  source_process_key: string;
  process_name: string;
  package_name: string;
  environment_name: string;
  successful_count: number;
  error_count: number;
  stopped_count: number;
  failed_count: number;
  executing_count: number;
  pending_count: number;
  suspended_count: number;
  resumed_count: number;
};

export type UseCaseDetail = UseCaseListItem & {
  source_key: string;
  updated_at: string;
  processes: UseCaseProcess[];
};

export type UseCaseResponse = {
  freshness: DataFreshness;
  use_case: UseCaseDetail;
};
