import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import EmailHistoryPage from "@/app/(protected)/email-history/page";
import EmailManagementPage from "@/app/(protected)/email-management/page";
import { getRpaDataProvider } from "@/lib/data/provider";
import type {
  EmailRecord,
  RpaDataProvider,
  UseCase,
} from "@/lib/data/types";

vi.mock("@/components/email-management-preview", () => ({
  EmailManagementPreview: ({
    recipients,
    reportTypes,
    useCases,
  }: {
    recipients: string[];
    reportTypes: string[];
    useCases: UseCase[];
  }) => (
    <p>
      Management data: {useCases.length}, {reportTypes.length},{" "}
      {recipients.length}
    </p>
  ),
}));

vi.mock("@/components/email-history-preview", () => ({
  EmailHistoryPreview: ({
    emailHistory,
  }: {
    emailHistory: EmailRecord[];
  }) => (
    <p>
      Email history data: {emailHistory.map((email) => email.id).join(", ")}
    </p>
  ),
}));

vi.mock("@/lib/data/provider", () => ({
  getRpaDataProvider: vi.fn(),
}));

const mockGetRpaDataProvider = vi.mocked(getRpaDataProvider);

beforeEach(() => {
  vi.clearAllMocks();
});

describe("Email route pages", () => {
  it("loads email-management data through the selected provider", async () => {
    const provider = {
      getUseCases: vi.fn().mockResolvedValue([{ id: "use-case-a" }]),
      getReportTypes: vi.fn().mockResolvedValue(["Performance Summary"]),
      getRecipients: vi.fn().mockResolvedValue(["preview@example.com"]),
    } as unknown as RpaDataProvider;
    mockGetRpaDataProvider.mockReturnValue(provider);

    render(await EmailManagementPage());

    expect(screen.getByText("Management data: 1, 1, 1")).toBeVisible();
    expect(provider.getUseCases).toHaveBeenCalledOnce();
    expect(provider.getReportTypes).toHaveBeenCalledOnce();
    expect(provider.getRecipients).toHaveBeenCalledOnce();
  });

  it("loads email history through the selected provider", async () => {
    const provider = {
      getEmailHistory: vi.fn().mockResolvedValue([{ id: "EM-5512" }]),
    } as unknown as RpaDataProvider;
    mockGetRpaDataProvider.mockReturnValue(provider);

    render(await EmailHistoryPage());

    expect(screen.getByText("Email history data: EM-5512")).toBeVisible();
    expect(provider.getEmailHistory).toHaveBeenCalledOnce();
  });
});
