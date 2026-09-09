import { fireEvent, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";

import { EmailHistoryPreview } from "@/components/email-history-preview";
import { EmailManagementPreview } from "@/components/email-management-preview";
import type { EmailRecord, UseCase } from "@/lib/data/types";

const toast = vi.hoisted(() => ({
  error: vi.fn(),
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

const emailHistory: EmailRecord[] = [
  {
    id: "EM-5512",
    useCase: "Use Case A",
    recipient: "finance.ops@company.com",
    subject: "Monthly Automation Report",
    date: "04 Sep 2026",
    status: "Sent",
  },
  {
    id: "EM-5509",
    useCase: "Use Case E",
    recipient: "compliance@company.com",
    subject: "Regulatory Submission Summary",
    date: "02 Sep 2026",
    status: "Failed",
  },
];

afterEach(() => {
  vi.clearAllMocks();
});

describe("EmailManagementPreview", () => {
  it("creates a schedule in local preview state without sending or saving", () => {
    render(
      <EmailManagementPreview
        recipients={["preview@example.com"]}
        reportTypes={["Performance Summary"]}
        useCases={[useCase]}
      />,
    );

    fireEvent.click(screen.getByRole("checkbox"));
    fireEvent.click(
      screen.getByRole("button", { name: "Simulate schedule" }),
    );

    expect(screen.getByText("Email schedules (4)")).toBeVisible();
    expect(toast.success).toHaveBeenCalledWith(
      "Schedule creation simulated",
      expect.objectContaining({
        description: expect.stringMatching(
          /Nothing was saved and no email was sent/i,
        ),
      }),
    );
  });

  it("removes a schedule from local preview state only", () => {
    render(
      <EmailManagementPreview
        recipients={["preview@example.com"]}
        reportTypes={["Performance Summary"]}
        useCases={[useCase]}
      />,
    );

    fireEvent.click(
      screen.getByRole("button", {
        name: "Remove preview schedule SCH-01",
      }),
    );

    expect(
      screen.queryByRole("button", {
        name: "Remove preview schedule SCH-01",
      }),
    ).not.toBeInTheDocument();
    expect(toast.success).toHaveBeenCalledWith(
      "Schedule removal simulated",
      expect.objectContaining({
        description: expect.stringMatching(/No saved schedule was changed/i),
      }),
    );
  });
});

describe("EmailHistoryPreview", () => {
  it("filters the provided history without mutating it", () => {
    render(<EmailHistoryPreview emailHistory={emailHistory} />);

    fireEvent.change(
      screen.getByRole("textbox", { name: "Search email history" }),
      { target: { value: "compliance" } },
    );

    expect(screen.queryByText("EM-5512")).not.toBeInTheDocument();
    expect(screen.getByText("EM-5509")).toBeVisible();
    expect(emailHistory).toHaveLength(2);
  });

  it("simulates resend without queuing or sending email", () => {
    render(<EmailHistoryPreview emailHistory={emailHistory} />);

    fireEvent.click(
      screen.getByRole("button", { name: "Preview resend for EM-5509" }),
    );

    expect(toast.success).toHaveBeenCalledWith(
      "Resend preview simulated",
      expect.objectContaining({
        description: expect.stringMatching(/was not queued or sent/i),
      }),
    );
  });
});
