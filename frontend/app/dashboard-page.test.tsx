import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import DashboardPage from "@/app/(protected)/dashboard/page";
import { getCurrentUser, getDashboardOverview } from "@/lib/api/server";
import type { DashboardOverviewResult } from "@/lib/api/server";
import {
  getRpaDataMode,
  getRpaDataProvider,
} from "@/lib/data/provider";
import type { RpaDataProvider } from "@/lib/data/types";

vi.mock("@/components/api-overview-dashboard", () => ({
  ApiOverviewDashboard: ({
    summary,
    trend,
    errors,
    canTriggerSync,
  }: {
    summary: { status: string; data?: { kpis: { total_use_cases: number } } };
    trend: { status: string; data?: { points: unknown[] } };
    errors: { status: string; data?: { groups: unknown[] } };
    canTriggerSync?: boolean;
  }) => (
    <div>
      <p>
        API values: {summary.data?.kpis.total_use_cases},{" "}
        {trend.data?.points.length}, {errors.data?.groups.length}
      </p>
      <p>Sync control: {canTriggerSync ? "shown" : "hidden"}</p>
    </div>
  ),
}));

vi.mock("@/components/overview-dashboard", () => ({
  OverviewDashboard: ({
    summary,
    performanceSeries,
    errorDistribution,
    automationTypeSplit,
    useCases,
  }: {
    summary: { totalUseCase: number };
    performanceSeries: unknown[];
    errorDistribution: unknown[];
    automationTypeSplit: unknown[];
    useCases: unknown[];
  }) => (
    <div>
      Provider values: {summary.totalUseCase}, {performanceSeries.length},{" "}
      {errorDistribution.length}, {automationTypeSplit.length}, {useCases.length}
    </div>
  ),
}));

vi.mock("@/lib/api/server", () => ({
  getCurrentUser: vi.fn(),
  getDashboardOverview: vi.fn(),
}));

vi.mock("@/lib/data/provider", () => ({
  getRpaDataMode: vi.fn(),
  getRpaDataProvider: vi.fn(),
}));

const mockGetCurrentUser = vi.mocked(getCurrentUser);
const mockGetDashboardOverview = vi.mocked(getDashboardOverview);

const mockGetRpaDataMode = vi.mocked(getRpaDataMode);
const mockGetRpaDataProvider = vi.mocked(getRpaDataProvider);

const viewer = {
  id: "user-1",
  employee_id: "10000001",
  display_name: "Viewer",
  role: "viewer" as const,
};

const admin = { ...viewer, id: "user-2", display_name: "Admin", role: "admin" as const };

const period = {
  from: "2026-08-07T00:00:00Z",
  to: "2026-09-06T00:00:00Z",
  timezone: "UTC",
};

const successOverview: DashboardOverviewResult = {
  status: "success",
  summary: {
    status: "success",
    data: {
      period,
      freshness: {
        status: "fresh",
        last_successful_refresh_at: "2026-09-06T00:00:00Z",
      },
      kpis: {
        total_use_cases: 42,
        active_use_cases: 39,
        execution_volume: 110,
        success_rate: 90,
        failed_executions: 11,
      },
    },
  },
  trend: {
    status: "success",
    data: {
      period,
      points: [
        {
          bucket: "2026-09-01T00:00:00Z",
          label: "Sep",
          success: 99,
          failure: 11,
        },
      ],
    },
  },
  errors: {
    status: "success",
    data: {
      period,
      groups: [{ code: "timeout", label: "Timeout", count: 11 }],
    },
  },
};

describe("DashboardPage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("uses only the three dashboard API resources in api mode", async () => {
    mockGetRpaDataMode.mockReturnValue("api");
    mockGetCurrentUser.mockResolvedValue(viewer);
    mockGetDashboardOverview.mockResolvedValue(successOverview);

    render(await DashboardPage());

    expect(screen.getByText("API values: 42, 1, 1")).toBeVisible();
    expect(mockGetDashboardOverview).toHaveBeenCalledOnce();
    expect(mockGetRpaDataProvider).not.toHaveBeenCalled();
  });

  it("withholds the sync control from a viewer", async () => {
    mockGetRpaDataMode.mockReturnValue("api");
    mockGetCurrentUser.mockResolvedValue(viewer);
    mockGetDashboardOverview.mockResolvedValue(successOverview);

    render(await DashboardPage());

    expect(screen.getByText("Sync control: hidden")).toBeVisible();
  });

  it("gives the sync control to an admin", async () => {
    mockGetRpaDataMode.mockReturnValue("api");
    mockGetCurrentUser.mockResolvedValue(admin);
    mockGetDashboardOverview.mockResolvedValue(successOverview);

    render(await DashboardPage());

    expect(screen.getByText("Sync control: shown")).toBeVisible();
  });

  it("preserves the complete FE-06 Overview in mock mode", async () => {
    mockGetRpaDataMode.mockReturnValue("mock");
    const provider = {
      getSummary: vi.fn().mockResolvedValue({
        totalUseCase: 42,
        successRate: 91,
        totalIssue: 8,
        unattended: 12,
        attended: 9,
        reportSent: 31,
      }),
      getPerformanceSeries: vi.fn().mockResolvedValue([
        { month: "Sep", success: 10, failed: 1, rate: 91 },
      ]),
      getErrorDistribution: vi
        .fn()
        .mockResolvedValue([{ name: "Error A", value: 8 }]),
      getAutomationTypeSplit: vi
        .fn()
        .mockResolvedValue([{ name: "Unattended", value: 12 }]),
      getUseCases: vi.fn().mockResolvedValue([
        {
          id: "use-case-a",
          name: "Use Case A",
          owner: "Operations",
          description: "Example",
          successRate: 91,
          status: "Running",
          automationType: "Unattended",
          totalProcess: 11,
          totalSuccess: 10,
          totalFailed: 1,
          issues: [],
          trend: [],
        },
      ]),
    } as unknown as RpaDataProvider;
    mockGetRpaDataProvider.mockReturnValue(provider);

    render(await DashboardPage());

    expect(screen.getByText("Provider values: 42, 1, 1, 1, 1")).toBeVisible();
    expect(provider.getSummary).toHaveBeenCalledOnce();
    expect(provider.getPerformanceSeries).toHaveBeenCalledOnce();
    expect(provider.getErrorDistribution).toHaveBeenCalledOnce();
    expect(provider.getAutomationTypeSplit).toHaveBeenCalledOnce();
    expect(provider.getUseCases).toHaveBeenCalledOnce();
    expect(mockGetDashboardOverview).not.toHaveBeenCalled();
  });
});
