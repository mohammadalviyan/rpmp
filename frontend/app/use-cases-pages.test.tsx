import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import UseCaseDetailPage from "@/app/(protected)/use-cases/[id]/page";
import UseCasesPage from "@/app/(protected)/use-cases/page";
import { UseCaseNotFound } from "@/components/use-case-detail";
import { getRpaDataProvider } from "@/lib/data/provider";
import type { RpaDataProvider, UseCase } from "@/lib/data/types";

vi.mock("@/components/use-cases-list", () => ({
  UseCasesList: ({ useCases }: { useCases: UseCase[] }) => (
    <p>Loaded list: {useCases.map((useCase) => useCase.name).join(", ")}</p>
  ),
}));

vi.mock("@/components/use-case-detail", async (importOriginal) => {
  const original =
    await importOriginal<typeof import("@/components/use-case-detail")>();

  return {
    ...original,
    UseCaseDetail: ({ useCase }: { useCase: UseCase }) => (
      <p>Loaded detail: {useCase.name}</p>
    ),
  };
});

vi.mock("@/lib/data/provider", () => ({
  getRpaDataProvider: vi.fn(),
}));

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

const mockGetRpaDataProvider = vi.mocked(getRpaDataProvider);

describe("Use Cases pages", () => {
  it("loads the list only through getUseCases", async () => {
    const getUseCases = vi.fn().mockResolvedValue([useCase]);
    const provider = { getUseCases } as unknown as RpaDataProvider;
    mockGetRpaDataProvider.mockReturnValue(provider);

    render(await UseCasesPage());

    expect(screen.getByText("Loaded list: Use Case A")).toBeVisible();
    expect(getUseCases).toHaveBeenCalledOnce();
  });

  it("loads one detail only through getUseCase", async () => {
    const getUseCase = vi.fn().mockResolvedValue(useCase);
    const provider = { getUseCase } as unknown as RpaDataProvider;
    mockGetRpaDataProvider.mockReturnValue(provider);

    render(
      await UseCaseDetailPage({
        params: Promise.resolve({ id: "use-case-a" }),
      }),
    );

    expect(screen.getByText("Loaded detail: Use Case A")).toBeVisible();
    expect(getUseCase).toHaveBeenCalledWith("use-case-a");
  });

  it("shows a clear not-found state for an unknown ID", async () => {
    const getUseCase = vi.fn().mockResolvedValue(undefined);
    const provider = { getUseCase } as unknown as RpaDataProvider;
    mockGetRpaDataProvider.mockReturnValue(provider);

    render(
      await UseCaseDetailPage({
        params: Promise.resolve({ id: "unknown" }),
      }),
    );

    expect(
      screen.getByRole("heading", { name: "This use case does not exist" }),
    ).toBeVisible();
    expect(
      screen.getByRole("link", { name: "Back to Use Cases" }),
    ).toHaveAttribute("href", "/use-cases");
    expect(getUseCase).toHaveBeenCalledWith("unknown");
  });
});

describe("UseCaseNotFound", () => {
  it("explains how to return to the list", () => {
    render(<UseCaseNotFound />);

    expect(screen.getByText("Use case not found")).toBeVisible();
  });
});
