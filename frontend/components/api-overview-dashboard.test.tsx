import { render, screen, within } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { ApiOverviewDashboard } from "@/components/api-overview-dashboard";
import type {
  DashboardErrors,
  DashboardSummary,
  ExecutionTrend,
} from "@/lib/api/types";
import type { DashboardResourceResult } from "@/lib/api/server";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ refresh: vi.fn() }),
}));

const period = {
  from: "2026-08-07T00:00:00Z",
  to: "2026-09-06T00:00:00Z",
  timezone: "UTC",
};

const summaryResult: DashboardResourceResult<DashboardSummary> = {
  status: "success",
  data: {
    period,
    freshness: {
      status: "fresh",
      last_successful_refresh_at: "2026-09-06T00:00:00Z",
    },
    kpis: {
      total_use_cases: 25,
      active_use_cases: 21,
      execution_volume: 250,
      success_rate: 92,
      failed_executions: 20,
    },
  },
};

const trendResult: DashboardResourceResult<ExecutionTrend> = {
  status: "success",
  data: {
    period,
    points: [
      {
        bucket: "2026-09-01T00:00:00Z",
        label: "Sep",
        success: 230,
        failure: 20,
      },
    ],
  },
};

const errorsResult: DashboardResourceResult<DashboardErrors> = {
  status: "success",
  data: {
    period,
    groups: [{ code: "timeout", label: "Timeout", count: 20 }],
  },
};

describe("ApiOverviewDashboard", () => {
  it("renders only the API-backed Overview data and stored-copy freshness", () => {
    render(
      <ApiOverviewDashboard
        summary={summaryResult}
        trend={trendResult}
        errors={errorsResult}
      />,
    );

    for (const [label, value] of [
      ["Total Use Cases", "25"],
      ["Active Use Cases", "21"],
      ["Execution Volume", "250"],
      ["Success Rate", "92%"],
      ["Failed Executions", "20"],
    ]) {
      const heading = screen.getByRole("heading", { name: label });
      expect(within(heading.closest("article")!).getByText(value)).toBeVisible();
    }

    expect(screen.getByText(/This is a stored RPMP copy/)).toHaveTextContent(
      "Status: fresh.",
    );
    expect(screen.getByText("Sep: 230 successful, 20 failed")).toBeInTheDocument();
    expect(
      screen.getByRole("list", { name: "Error distribution by category" }),
    ).toHaveTextContent("Timeout");

    expect(screen.queryByText("Unattended Processes")).not.toBeInTheDocument();
    expect(screen.queryByText("Attended Processes")).not.toBeInTheDocument();
    expect(screen.queryByText("Reports Sent")).not.toBeInTheDocument();
    expect(screen.queryByText("Automation Type")).not.toBeInTheDocument();
    expect(screen.queryByText("Top Use Cases")).not.toBeInTheDocument();
  });

  it("shows independent resource failures without mock fallbacks", () => {
    render(
      <ApiOverviewDashboard
        summary={{ status: "source_unavailable" }}
        trend={{ status: "unexpected_error" }}
        errors={{ status: "source_unavailable" }}
      />,
    );

    expect(
      screen.getByText(
        "Dashboard summary is unavailable because the stored RPMP data source could not be read.",
      ),
    ).toBeVisible();
    expect(screen.getByText("Execution trend could not be loaded. Try again later.")).toBeVisible();
    expect(
      screen.getByText(
        "Top errors is unavailable because the stored RPMP data source could not be read.",
      ),
    ).toBeVisible();
  });

  it("keeps the sync control out of the page for a viewer", () => {
    render(
      <ApiOverviewDashboard
        summary={summaryResult}
        trend={trendResult}
        errors={errorsResult}
      />,
    );

    expect(
      screen.queryByRole("button", { name: "Sync now" }),
    ).not.toBeInTheDocument();
    expect(screen.getByText(/This is a stored RPMP copy/)).toBeVisible();
  });

  it("offers the sync control to an admin", () => {
    render(
      <ApiOverviewDashboard
        canTriggerSync
        summary={summaryResult}
        trend={trendResult}
        errors={errorsResult}
      />,
    );

    expect(screen.getByRole("button", { name: "Sync now" })).toBeEnabled();
  });
});
