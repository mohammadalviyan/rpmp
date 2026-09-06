import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import DashboardPage from "@/app/(protected)/dashboard/page";
import { getDashboardOverview } from "@/lib/api/server";

const redirect = vi.fn();

vi.mock("next/navigation", () => ({
  redirect: (destination: string) => redirect(destination),
}));

vi.mock("@/lib/api/server", () => ({
  getDashboardOverview: vi.fn(),
}));

const mockGetDashboardOverview = vi.mocked(getDashboardOverview);

const summary = {
  period: {
    from: "2026-08-07T00:00:00Z",
    to: "2026-09-06T00:00:00Z",
    timezone: "UTC",
  },
  freshness: {
    status: "fresh",
    last_successful_refresh_at: "2026-09-06T00:00:00Z",
  },
  kpis: {
    total_use_cases: 12,
    active_use_cases: 9,
    execution_volume: 100,
    success_rate: 94,
    failed_executions: 6,
  },
};

describe("DashboardPage", () => {
  beforeEach(() => {
    redirect.mockReset();
    mockGetDashboardOverview.mockReset();
    redirect.mockImplementation(() => {
      throw new Error("NEXT_REDIRECT");
    });
  });

  it("redirects when any dashboard request rejects the session", async () => {
    mockGetDashboardOverview.mockResolvedValue({
      status: "unauthenticated",
    });

    await expect(DashboardPage()).rejects.toThrow("NEXT_REDIRECT");
    expect(redirect).toHaveBeenCalledWith("/login");
  });

  it("renders five KPIs and accessible chart values from API responses", async () => {
    mockGetDashboardOverview.mockResolvedValue({
      status: "success",
      summary: { status: "success", data: summary },
      trend: {
        status: "success",
        data: {
          period: summary.period,
          points: [
            {
              bucket: "2026-08-01T00:00:00Z",
              label: "Aug",
              success: 94,
              failure: 6,
            },
          ],
        },
      },
      errors: {
        status: "success",
        data: {
          period: summary.period,
          groups: [
            { code: "faulted", label: "Faulted", count: 4 },
            { code: "stopped", label: "Stopped", count: 2 },
          ],
        },
      },
    });

    render(await DashboardPage());

    expect(
      screen.getByRole("region", {
        name: "Dashboard key performance indicators",
      }),
    ).toBeVisible();
    expect(screen.getAllByRole("article")).toHaveLength(5);
    expect(
      screen.getByRole("img", { name: /Execution trend with 1 time points/ }),
    ).toBeVisible();
    expect(screen.getByText(/Aug: 94 successful, 6 failed/)).toBeInTheDocument();
    expect(screen.getByText("Faulted")).toBeVisible();
    expect(screen.getByText("4")).toBeVisible();
  });

  it("shows independent unavailable states without hiding successful panels", async () => {
    mockGetDashboardOverview.mockResolvedValue({
      status: "success",
      summary: { status: "source_unavailable" },
      trend: {
        status: "success",
        data: { period: summary.period, points: [] },
      },
      errors: { status: "unexpected_error" },
    });

    render(await DashboardPage());

    expect(
      screen.getByRole("heading", { name: "Dashboard data unavailable" }),
    ).toBeVisible();
    expect(
      screen.getByText(/Current dashboard data is temporarily unavailable/),
    ).toBeVisible();
    expect(screen.queryByText("Total Use Cases")).not.toBeInTheDocument();
    expect(screen.queryByText("12")).not.toBeInTheDocument();
    expect(screen.queryByText(/Period:/)).not.toBeInTheDocument();
    expect(screen.queryByText(/Last refreshed:/)).not.toBeInTheDocument();
    expect(
      screen.getByText("No execution trend data is available for this period."),
    ).toBeVisible();
    expect(
      screen.getByText(
        "We could not load error distribution. Try again later.",
      ),
    ).toBeVisible();
  });
});
