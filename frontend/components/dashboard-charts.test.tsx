import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { ErrorDistribution } from "@/components/error-distribution";
import { ExecutionTrendChart } from "@/components/execution-trend-chart";

const period = {
  from: "2026-08-07T00:00:00Z",
  to: "2026-09-06T00:00:00Z",
  timezone: "UTC",
};

describe("dashboard charts", () => {
  it("provides text values for every execution trend point", () => {
    render(
      <ExecutionTrendChart
        trend={{
          period,
          points: [
            {
              bucket: "2026-08-01T00:00:00Z",
              label: "Aug",
              success: 94,
              failure: 6,
            },
            {
              bucket: "2026-09-01T00:00:00Z",
              label: "Sep",
              success: 0,
              failure: 0,
            },
          ],
        }}
      />,
    );

    expect(
      screen.getByRole("img", { name: /Execution trend with 2 time points/ }),
    ).toBeVisible();
    expect(screen.getByText("Aug: 94 successful, 6 failed")).toBeInTheDocument();
    expect(screen.getByText("Sep: 0 successful, 0 failed")).toBeInTheDocument();
  });

  it("shows error labels and counts without relying on color", () => {
    render(
      <ErrorDistribution
        errors={{
          period,
          groups: [
            { code: "faulted", label: "Faulted", count: 4 },
            { code: "stopped", label: "Stopped", count: 2 },
          ],
        }}
      />,
    );

    const distribution = screen.getByRole("list", {
      name: "Error distribution by category",
    });
    expect(distribution).toHaveTextContent("Faulted");
    expect(distribution).toHaveTextContent("4");
    expect(distribution).toHaveTextContent("67%");
    expect(distribution).toHaveTextContent("Stopped");
    expect(distribution).toHaveTextContent("2");
    expect(distribution).toHaveTextContent("33%");
    expect(
      screen.getByRole("img", {
        name: "Error distribution donut chart. 6 failed executions across 2 categories.",
      }),
    ).toBeVisible();
  });

  it("shows explicit empty states", () => {
    const { rerender } = render(
      <ExecutionTrendChart trend={{ period, points: [] }} />,
    );
    expect(
      screen.getByText("No execution trend data is available for this period."),
    ).toBeVisible();

    rerender(
      <ErrorDistribution
        errors={{
          period,
          groups: [{ code: "faulted", label: "Faulted", count: 0 }],
        }}
      />,
    );
    expect(
      screen.getByText("No errors were recorded for this period."),
    ).toBeVisible();
  });
});
