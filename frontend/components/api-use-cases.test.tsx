import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { ApiUseCaseDetail } from "@/components/api-use-case-detail";
import { ApiUseCasesList } from "@/components/api-use-cases-list";
import type {
  UseCaseListItem,
  UseCaseResponse,
} from "@/lib/api/types";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

const listItems: UseCaseListItem[] = [
  {
    id: "invoice",
    name: "Invoice Processing",
    status: "active",
    process_count: 2,
    successful_count: 100,
    error_count: 3,
    stopped_count: 1,
    failed_count: 4,
    execution_volume: 104,
    success_rate: 96.15,
    environments: ["PROD"],
  },
  {
    id: "empty-volume",
    name: "New Use Case",
    status: "inactive",
    process_count: 0,
    successful_count: 0,
    error_count: 0,
    stopped_count: 0,
    failed_count: 0,
    execution_volume: 0,
    success_rate: null,
    environments: [],
  },
];

const detailResponse: UseCaseResponse = {
  freshness: {
    status: "fresh",
    last_successful_refresh_at: "2026-09-11T02:00:00Z",
  },
  use_case: {
    ...listItems[0],
    source_key: "Finance/Invoice.Process",
    updated_at: "2026-09-11T02:00:00Z",
    processes: [
      {
        source_process_key: "42",
        process_name: "Invoice",
        package_name: "Invoice.1.0.0",
        environment_name: "PROD",
        successful_count: 50,
        error_count: 1,
        stopped_count: 0,
        failed_count: 1,
        executing_count: 0,
        pending_count: 0,
        suspended_count: 0,
        resumed_count: 0,
      },
    ],
  },
};

describe("ApiUseCasesList", () => {
  it("renders API fields, null success rate, and no mock-only controls", () => {
    render(<ApiUseCasesList useCases={listItems} />);

    expect(screen.getByText("Invoice Processing")).toBeVisible();
    expect(screen.getByText("96.15%")).toBeVisible();
    expect(screen.getByText("Not available")).toBeVisible();
    expect(screen.getByText("No environment recorded")).toBeVisible();
    expect(screen.queryByRole("tab", { name: "Attended" })).not.toBeInTheDocument();
    expect(screen.queryByText("Automation Type")).not.toBeInTheDocument();
  });

  it("searches API rows and explains empty snapshots", async () => {
    const user = userEvent.setup();
    const { rerender } = render(<ApiUseCasesList useCases={listItems} />);

    await user.type(
      screen.getByRole("textbox", { name: "Search use cases" }),
      "invoice",
    );

    expect(screen.getByText("Invoice Processing")).toBeVisible();
    expect(screen.queryByText("New Use Case")).not.toBeInTheDocument();

    rerender(<ApiUseCasesList useCases={[]} />);
    expect(
      screen.getByText("No use cases are available in the latest snapshot."),
    ).toBeVisible();
  });
});

describe("ApiUseCaseDetail", () => {
  it("renders identity, four KPIs, process rows, and snapshot freshness", () => {
    render(<ApiUseCaseDetail data={detailResponse} />);

    expect(
      screen.getByRole("heading", { name: "Invoice Processing" }),
    ).toBeVisible();
    expect(screen.getByText("Source key: Finance/Invoice.Process")).toBeVisible();
    expect(
      screen.getByLabelText("Use case key performance indicators").children,
    ).toHaveLength(4);
    expect(screen.getByText("Invoice.1.0.0")).toBeVisible();
    expect(
      screen.getByText(/Last successful refresh: Sep 11, 2026, 2:00 AM UTC/),
    ).toBeVisible();
  });

  it("shows no fake weekly or issue history and handles a null success rate", () => {
    render(
      <ApiUseCaseDetail
        data={{
          ...detailResponse,
          use_case: {
            ...detailResponse.use_case,
            success_rate: null,
          },
        }}
      />,
    );

    expect(screen.getByText("Not available")).toBeVisible();
    expect(
      screen.getAllByText("This snapshot has no execution history."),
    ).toHaveLength(2);
    expect(screen.queryByText("Successful vs failed transactions per day")).not.toBeInTheDocument();
  });
});
