import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { UseCasesList } from "@/components/use-cases-list";
import type { UseCase } from "@/lib/data/types";

const push = vi.fn();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push }),
}));

const useCases: UseCase[] = [
  {
    id: "invoice-processing",
    name: "Invoice Processing",
    owner: "Finance Operations",
    description: "Invoice processing",
    successRate: 94,
    status: "Running",
    automationType: "Unattended",
    totalProcess: 100,
    totalSuccess: 94,
    totalFailed: 6,
    issues: [
      {
        id: "issue-a",
        name: "Timeout",
        description: "A timeout occurred",
        occurredAt: "03 Sep 2026, 09:14",
        status: "Open",
        errorType: "Error A",
      },
    ],
    trend: [],
  },
  {
    id: "ticket-assistant",
    name: "Ticket Assistant",
    owner: "Customer Service",
    description: "Ticket assistance",
    successRate: 88,
    status: "Warning",
    automationType: "Attended",
    totalProcess: 80,
    totalSuccess: 70,
    totalFailed: 10,
    issues: [],
    trend: [],
  },
];

describe("UseCasesList", () => {
  it("searches by use case name and owner", async () => {
    const user = userEvent.setup();
    render(<UseCasesList useCases={useCases} />);

    await user.type(screen.getByRole("textbox", { name: "Search use cases" }), "finance");

    expect(screen.getByText("Invoice Processing")).toBeVisible();
    expect(screen.queryByText("Ticket Assistant")).not.toBeInTheDocument();
  });

  it("filters the copied table with automation type tabs", async () => {
    const user = userEvent.setup();
    render(<UseCasesList useCases={useCases} />);

    await user.click(screen.getByRole("tab", { name: "Attended" }));

    expect(screen.getByText("Ticket Assistant")).toBeVisible();
    expect(screen.queryByText("Invoice Processing")).not.toBeInTheDocument();
  });

  it("supports row and detail-link navigation", async () => {
    const user = userEvent.setup();
    render(<UseCasesList useCases={useCases} />);

    expect(
      screen.getByRole("link", { name: "View Invoice Processing details" }),
    ).toHaveAttribute("href", "/use-cases/invoice-processing");

    await user.click(screen.getByLabelText("Open Invoice Processing"));

    expect(push).toHaveBeenCalledWith("/use-cases/invoice-processing");
  });
});
