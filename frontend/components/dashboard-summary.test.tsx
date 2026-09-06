import { render, screen, within } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { DashboardSummary } from "@/components/dashboard-summary";
import type { DashboardSummary as DashboardSummaryData } from "@/lib/api/types";

const summary: DashboardSummaryData = {
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

describe("DashboardSummary", () => {
  it("renders all five KPI values from the summary response", () => {
    render(<DashboardSummary summary={summary} />);

    const expectedKpis = [
      ["Total Use Cases", "12"],
      ["Active Use Cases", "9"],
      ["Execution Volume", "100"],
      ["Success Rate", "94%"],
      ["Failed Executions", "6"],
    ];

    for (const [label, value] of expectedKpis) {
      const heading = screen.getByRole("heading", { name: label });
      expect(within(heading.closest("article")!).getByText(value)).toBeVisible();
    }
  });

  it("shows the period, refresh time, timezone, and freshness as text", () => {
    render(<DashboardSummary summary={summary} />);

    expect(screen.getByText(/Period: Aug 7, 2026/)).toHaveTextContent(
      "Sep 6, 2026",
    );
    expect(screen.getByText(/Period:/)).toHaveTextContent("(UTC)");
    expect(screen.getByText(/Last refreshed: Sep 6, 2026/)).toHaveTextContent(
      "Freshness: fresh.",
    );
  });

  it("renders empty response values as current data and null success rate as N/A", () => {
    render(
      <DashboardSummary
        summary={{
          ...summary,
          kpis: {
            total_use_cases: 0,
            active_use_cases: 0,
            execution_volume: 0,
            success_rate: null,
            failed_executions: 0,
          },
        }}
      />,
    );

    const successRateHeading = screen.getByRole("heading", {
      name: "Success Rate",
    });
    expect(
      within(successRateHeading.closest("article")!).getByText("N/A"),
    ).toBeVisible();
    expect(screen.queryByText("0%")).not.toBeInTheDocument();
    expect(
      screen.getByText("There were no executions in this period."),
    ).toBeVisible();
    expect(screen.getAllByText("0")).toHaveLength(4);
  });
});
