import { act, fireEvent, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";

import { GenerateReportPreview } from "@/components/generate-report-preview";
import { ReportHistoryPreview } from "@/components/report-history-preview";
import type { ReportRecord, UseCase } from "@/lib/data/types";

const toast = vi.hoisted(() => ({
  error: vi.fn(),
  info: vi.fn(),
  success: vi.fn(),
}));

vi.mock("sonner", () => ({ toast }));

const useCase: UseCase = {
  id: "use-case-a",
  name: "Use Case A",
  owner: "Finance Operations",
  description: "Invoice processing",
  successRate: 94,
  status: "Running",
  automationType: "Unattended",
  totalProcess: 100,
  totalSuccess: 94,
  totalFailed: 6,
  issues: [],
  trend: [],
};

const report: ReportRecord = {
  id: "RPT-2091",
  useCase: "Use Case A",
  type: "Performance Summary",
  period: "Aug 2026",
  generatedAt: "04 Sep 2026, 06:00",
  size: "1.2 MB",
  status: "Ready",
};

afterEach(() => {
  vi.clearAllMocks();
  vi.useRealTimers();
});

describe("GenerateReportPreview", () => {
  it("labels the page and draft action as preview-only", () => {
    render(
      <GenerateReportPreview
        recipients={["preview@example.com"]}
        reportPeriods={["Today", "Last 7 days", "Last 30 days", "Aug 2026"]}
        reportTypes={["Performance Summary"]}
        useCases={[useCase]}
      />,
    );

    expect(screen.getByText("Preview only")).toBeVisible();
    fireEvent.click(
      screen.getByRole("button", { name: "Simulate draft save" }),
    );

    expect(toast.info).toHaveBeenCalledWith(
      "Draft save simulated",
      expect.objectContaining({
        description: expect.stringMatching(
          /No draft was saved and report history was not changed/i,
        ),
      }),
    );
  });

  it("simulates generation without sending email or changing history", () => {
    vi.useFakeTimers();
    render(
      <GenerateReportPreview
        recipients={["preview@example.com"]}
        reportPeriods={["Today", "Last 7 days", "Last 30 days", "Aug 2026"]}
        reportTypes={["Performance Summary"]}
        useCases={[useCase]}
      />,
    );

    fireEvent.click(screen.getByRole("button", { name: "Simulate report" }));
    act(() => vi.advanceTimersByTime(800));

    expect(toast.success).toHaveBeenCalledWith(
      "Report generation simulated",
      expect.objectContaining({
        description: expect.stringMatching(
          /No email was sent and no history was changed/i,
        ),
      }),
    );
  });
});

describe("ReportHistoryPreview", () => {
  it("keeps download preview-only and does not render a download link", () => {
    const { container } = render(
      <ReportHistoryPreview reportHistory={[report]} />,
    );

    expect(screen.getByText("Demo data")).toBeVisible();
    expect(container.querySelector("a[download]")).not.toBeInTheDocument();

    fireEvent.click(
      screen.getByRole("button", {
        name: "Preview download for RPT-2091",
      }),
    );

    expect(toast.info).toHaveBeenCalledWith(
      "Download preview simulated",
      expect.objectContaining({
        description: expect.stringMatching(
          /was not downloaded.*no file was created/i,
        ),
      }),
    );
  });

  it("filters the provided report history without mutating it", () => {
    const history = [
      report,
      {
        ...report,
        id: "RPT-2090",
        useCase: "Use Case C",
        type: "Issue Analysis",
      },
    ];
    render(<ReportHistoryPreview reportHistory={history} />);

    fireEvent.change(screen.getByRole("textbox", { name: "Search reports" }), {
      target: { value: "Issue Analysis" },
    });

    expect(screen.queryByText("RPT-2091")).not.toBeInTheDocument();
    expect(screen.getByText("RPT-2090")).toBeVisible();
    expect(history).toHaveLength(2);
  });
});
