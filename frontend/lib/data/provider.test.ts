import { describe, expect, it } from "vitest";

import { apiDataProvider } from "@/lib/data/api-provider";
import { mockDataProvider } from "@/lib/data/mock-provider";
import {
  getRpaDataMode,
  getRpaDataProvider,
} from "@/lib/data/provider";

describe("RPMP data providers", () => {
  it("defaults unknown and missing modes to mock", () => {
    expect(getRpaDataMode(undefined)).toBe("mock");
    expect(getRpaDataMode("preview")).toBe("mock");
    expect(getRpaDataProvider("mock")).toBe(mockDataProvider);
  });

  it("selects the typed API stub only in api mode", () => {
    expect(getRpaDataMode("api")).toBe("api");
    expect(getRpaDataProvider("api")).toBe(apiDataProvider);
  });

  it("exposes the complete sample provider shape", async () => {
    const [
      useCases,
      summary,
      performance,
      errors,
      automationTypes,
      reports,
      emails,
      recipients,
      reportTypes,
      reportPeriods,
    ] = await Promise.all([
      mockDataProvider.getUseCases(),
      mockDataProvider.getSummary(),
      mockDataProvider.getPerformanceSeries(),
      mockDataProvider.getErrorDistribution(),
      mockDataProvider.getAutomationTypeSplit(),
      mockDataProvider.getReportHistory(),
      mockDataProvider.getEmailHistory(),
      mockDataProvider.getRecipients(),
      mockDataProvider.getReportTypes(),
      mockDataProvider.getReportPeriods(),
    ]);

    expect(useCases).toHaveLength(5);
    expect(await mockDataProvider.getUseCase("use-case-a")).toMatchObject({
      owner: "Finance Operations",
      successRate: 94,
    });
    expect(summary).toEqual({
      totalUseCase: 25,
      successRate: 85,
      totalIssue: 13,
      unattended: 14,
      attended: 11,
      reportSent: 20,
    });
    expect(performance).toHaveLength(7);
    expect(errors).toHaveLength(2);
    expect(automationTypes).toHaveLength(2);
    expect(reports).toHaveLength(5);
    expect(emails).toHaveLength(5);
    expect(recipients).toHaveLength(6);
    expect(reportTypes).toHaveLength(4);
    expect(reportPeriods).toHaveLength(5);
  });
});
