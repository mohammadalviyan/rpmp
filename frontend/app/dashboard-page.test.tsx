import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import DashboardPage from "@/app/(protected)/dashboard/page";
import { getDashboardSummary } from "@/lib/api/server";

const redirect = vi.fn();

vi.mock("next/navigation", () => ({
  redirect: (destination: string) => redirect(destination),
}));

vi.mock("@/lib/api/server", () => ({
  getDashboardSummary: vi.fn(),
}));

const mockGetDashboardSummary = vi.mocked(getDashboardSummary);

describe("DashboardPage", () => {
  beforeEach(() => {
    redirect.mockReset();
    mockGetDashboardSummary.mockReset();
    redirect.mockImplementation(() => {
      throw new Error("NEXT_REDIRECT");
    });
  });

  it("redirects a rejected summary request to login", async () => {
    mockGetDashboardSummary.mockResolvedValue({
      status: "unauthenticated",
    });

    await expect(DashboardPage()).rejects.toThrow("NEXT_REDIRECT");
    expect(redirect).toHaveBeenCalledWith("/login");
  });

  it("shows source-unavailable copy without KPI values or current timestamps", async () => {
    mockGetDashboardSummary.mockResolvedValue({
      status: "source_unavailable",
    });

    render(await DashboardPage());

    expect(
      screen.getByRole("heading", { name: "Dashboard data unavailable" }),
    ).toBeVisible();
    expect(
      screen.getByText(/Current dashboard data is temporarily unavailable/),
    ).toBeVisible();
    expect(screen.queryByText("Total Use Cases")).not.toBeInTheDocument();
    expect(screen.queryByText("12")).not.toBeInTheDocument();
    expect(screen.queryByText(/Period:/)).not.toBeInTheDocument();
    expect(screen.queryByText(/Last refreshed:/)).not.toBeInTheDocument();
  });

  it("shows safe copy for unexpected failures", async () => {
    mockGetDashboardSummary.mockResolvedValue({
      status: "unexpected_error",
    });

    render(await DashboardPage());

    expect(
      screen.getByText("We could not load current dashboard data. Try again later."),
    ).toBeVisible();
    expect(screen.queryByText("Execution Volume")).not.toBeInTheDocument();
  });
});
