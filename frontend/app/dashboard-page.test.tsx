import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import DashboardPage from "@/app/(protected)/dashboard/page";
import { getRpaDataProvider } from "@/lib/data/provider";
import type { RpaDataProvider } from "@/lib/data/types";

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

vi.mock("@/lib/data/provider", () => ({
  getRpaDataProvider: vi.fn(),
}));

const mockGetRpaDataProvider = vi.mocked(getRpaDataProvider);

describe("DashboardPage", () => {
  it("loads every overview section through the selected provider", async () => {
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
  });
});
