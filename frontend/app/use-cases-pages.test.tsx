import { render, screen } from "@testing-library/react";
import { redirect } from "next/navigation";
import { beforeEach, describe, expect, it, vi } from "vitest";

import UseCaseDetailPage from "@/app/(protected)/use-cases/[id]/page";
import UseCasesPage from "@/app/(protected)/use-cases/page";
import { UseCaseNotFound } from "@/components/use-case-detail";
import { getUseCase, getUseCases } from "@/lib/api/server";
import type { UseCaseResponse, UseCasesResponse } from "@/lib/api/types";
import {
  getRpaDataMode,
  getRpaDataProvider,
} from "@/lib/data/provider";
import type { RpaDataProvider, UseCase } from "@/lib/data/types";

vi.mock("next/navigation", () => ({
  redirect: vi.fn(),
}));

vi.mock("@/components/api-use-cases-list", () => ({
  ApiUseCasesList: ({ useCases }: { useCases: UseCasesResponse["items"] }) => (
    <p>API list: {useCases.map((useCase) => useCase.name).join(", ")}</p>
  ),
}));

vi.mock("@/components/api-use-case-detail", () => ({
  ApiUseCaseDetail: ({ data }: { data: UseCaseResponse }) => (
    <p>API detail: {data.use_case.name}</p>
  ),
}));

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

vi.mock("@/lib/api/server", () => ({
  getUseCase: vi.fn(),
  getUseCases: vi.fn(),
}));

vi.mock("@/lib/data/provider", () => ({
  getRpaDataMode: vi.fn(),
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
const mockGetRpaDataMode = vi.mocked(getRpaDataMode);
const mockGetUseCase = vi.mocked(getUseCase);
const mockGetUseCases = vi.mocked(getUseCases);
const mockRedirect = vi.mocked(redirect);

describe("Use Cases pages", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("preserves the mock list path", async () => {
    mockGetRpaDataMode.mockReturnValue("mock");
    const getProviderUseCases = vi.fn().mockResolvedValue([useCase]);
    const provider = {
      getUseCases: getProviderUseCases,
    } as unknown as RpaDataProvider;
    mockGetRpaDataProvider.mockReturnValue(provider);

    render(await UseCasesPage());

    expect(screen.getByText("Loaded list: Use Case A")).toBeVisible();
    expect(getProviderUseCases).toHaveBeenCalledOnce();
    expect(mockGetUseCases).not.toHaveBeenCalled();
  });

  it("preserves the mock detail path", async () => {
    mockGetRpaDataMode.mockReturnValue("mock");
    const getProviderUseCase = vi.fn().mockResolvedValue(useCase);
    const provider = {
      getUseCase: getProviderUseCase,
    } as unknown as RpaDataProvider;
    mockGetRpaDataProvider.mockReturnValue(provider);

    render(
      await UseCaseDetailPage({
        params: Promise.resolve({ id: "use-case-a" }),
      }),
    );

    expect(screen.getByText("Loaded detail: Use Case A")).toBeVisible();
    expect(getProviderUseCase).toHaveBeenCalledWith("use-case-a");
    expect(mockGetUseCase).not.toHaveBeenCalled();
  });

  it("shows a clear not-found state for an unknown mock ID", async () => {
    mockGetRpaDataMode.mockReturnValue("mock");
    const getProviderUseCase = vi.fn().mockResolvedValue(undefined);
    const provider = {
      getUseCase: getProviderUseCase,
    } as unknown as RpaDataProvider;
    mockGetRpaDataProvider.mockReturnValue(provider);

    render(
      await UseCaseDetailPage({
        params: Promise.resolve({ id: "unknown" }),
      }),
    );

    expect(
      screen.getByRole("heading", { name: "This use case does not exist" }),
    ).toBeVisible();
    expect(getProviderUseCase).toHaveBeenCalledWith("unknown");
  });

  it("uses only the RPMP list endpoint in api mode", async () => {
    mockGetRpaDataMode.mockReturnValue("api");
    mockGetUseCases.mockResolvedValue({
      status: "success",
      data: {
        freshness: {
          status: "fresh",
          last_successful_refresh_at: "2026-09-11T02:00:00Z",
        },
        items: [
          {
            id: "api-use-case",
            name: "API Use Case",
            status: "active",
            process_count: 1,
            successful_count: 10,
            error_count: 1,
            stopped_count: 0,
            failed_count: 1,
            execution_volume: 11,
            success_rate: 90.91,
            environments: ["PROD"],
          },
        ],
      },
    });

    render(await UseCasesPage());

    expect(screen.getByText("API list: API Use Case")).toBeVisible();
    expect(mockGetUseCases).toHaveBeenCalledOnce();
    expect(mockGetRpaDataProvider).not.toHaveBeenCalled();
  });

  it("uses only the RPMP detail endpoint in api mode", async () => {
    mockGetRpaDataMode.mockReturnValue("api");
    mockGetUseCase.mockResolvedValue({
      status: "success",
      data: {
        freshness: {
          status: "fresh",
          last_successful_refresh_at: "2026-09-11T02:00:00Z",
        },
        use_case: {
          id: "api-use-case",
          name: "API Use Case",
          status: "active",
          source_key: "Finance/API.UseCase",
          updated_at: "2026-09-11T02:00:00Z",
          process_count: 0,
          successful_count: 0,
          error_count: 0,
          stopped_count: 0,
          failed_count: 0,
          execution_volume: 0,
          success_rate: null,
          environments: [],
          processes: [],
        },
      },
    });

    render(
      await UseCaseDetailPage({
        params: Promise.resolve({ id: "api-use-case" }),
      }),
    );

    expect(screen.getByText("API detail: API Use Case")).toBeVisible();
    expect(mockGetUseCase).toHaveBeenCalledWith("api-use-case");
    expect(mockGetRpaDataProvider).not.toHaveBeenCalled();
  });

  it("maps an API detail 404 to the existing not-found UI", async () => {
    mockGetRpaDataMode.mockReturnValue("api");
    mockGetUseCase.mockResolvedValue({ status: "not_found" });

    render(
      await UseCaseDetailPage({
        params: Promise.resolve({ id: "unknown" }),
      }),
    );

    expect(
      screen.getByRole("heading", { name: "This use case does not exist" }),
    ).toBeVisible();
  });

  it("redirects an unauthenticated API list using the Overview pattern", async () => {
    mockGetRpaDataMode.mockReturnValue("api");
    mockGetUseCases.mockResolvedValue({ status: "unauthenticated" });

    await UseCasesPage();

    expect(mockRedirect).toHaveBeenCalledWith("/login");
  });
});

describe("UseCaseNotFound", () => {
  it("explains how to return to the list", () => {
    render(<UseCaseNotFound />);

    expect(screen.getByText("Use case not found")).toBeVisible();
  });
});
