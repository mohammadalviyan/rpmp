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
