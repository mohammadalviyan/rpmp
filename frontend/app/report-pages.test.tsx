import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import GenerateReportPage from "@/app/(protected)/generate-report/page";
import ReportHistoryPage from "@/app/(protected)/report-history/page";
import { getRpaDataProvider } from "@/lib/data/provider";
import type {
  ReportRecord,
  RpaDataProvider,
  UseCase,
} from "@/lib/data/types";

vi.mock("@/components/generate-report-preview", () => ({
  GenerateReportPreview: ({
    recipients,
    reportPeriods,
    reportTypes,
    useCases,
  }: {
    recipients: string[];
    reportPeriods: string[];
    reportTypes: string[];
    useCases: UseCase[];
  }) => (
    <p>
      Generate data: {useCases.length}, {reportTypes.length},{" "}
      {reportPeriods.length}, {recipients.length}
    </p>
  ),
}));

vi.mock("@/components/report-history-preview", () => ({
  ReportHistoryPreview: ({
    reportHistory,
  }: {
    reportHistory: ReportRecord[];
  }) => <p>History data: {reportHistory.map((report) => report.id).join(", ")}</p>,
}));

vi.mock("@/lib/data/provider", () => ({
  getRpaDataProvider: vi.fn(),
}));

const mockGetRpaDataProvider = vi.mocked(getRpaDataProvider);

describe("Report route pages", () => {
  it("loads generate-report data through the selected provider", async () => {
    const provider = {
      getUseCases: vi.fn().mockResolvedValue([{ id: "use-case-a" }]),
      getReportTypes: vi.fn().mockResolvedValue(["Performance Summary"]),
      getReportPeriods: vi.fn().mockResolvedValue(["Last 30 days"]),
      getRecipients: vi.fn().mockResolvedValue(["preview@example.com"]),
    } as unknown as RpaDataProvider;
    mockGetRpaDataProvider.mockReturnValue(provider);

    render(await GenerateReportPage());

    expect(screen.getByText("Generate data: 1, 1, 1, 1")).toBeVisible();
    expect(provider.getUseCases).toHaveBeenCalledOnce();
    expect(provider.getReportTypes).toHaveBeenCalledOnce();
    expect(provider.getReportPeriods).toHaveBeenCalledOnce();
    expect(provider.getRecipients).toHaveBeenCalledOnce();
  });

  it("loads report history through the selected provider", async () => {
    const provider = {
      getReportHistory: vi.fn().mockResolvedValue([{ id: "RPT-2091" }]),
    } as unknown as RpaDataProvider;
    mockGetRpaDataProvider.mockReturnValue(provider);

    render(await ReportHistoryPage());

    expect(screen.getByText("History data: RPT-2091")).toBeVisible();
    expect(provider.getReportHistory).toHaveBeenCalledOnce();
  });
});
