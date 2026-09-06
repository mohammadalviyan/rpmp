import "server-only";

import { cookies } from "next/headers";

import type {
  AuthResponse,
  DashboardErrors,
  DashboardSummary,
  ExecutionTrend,
  User,
} from "@/lib/api/types";

const apiOrigin = (
  process.env.RPMP_API_ORIGIN ?? "http://127.0.0.1:8080"
).replace(/\/$/, "");

export type DashboardSummaryResult =
  | { status: "success"; summary: DashboardSummary }
  | { status: "unauthenticated" }
  | { status: "source_unavailable" }
  | { status: "unexpected_error" };

export type DashboardResourceResult<T> =
  | { status: "success"; data: T }
  | { status: "unauthenticated" }
  | { status: "source_unavailable" }
  | { status: "unexpected_error" };

export type DashboardOverviewResult =
  | {
      status: "success";
      summary: DashboardResourceResult<DashboardSummary>;
      trend: DashboardResourceResult<ExecutionTrend>;
      errors: DashboardResourceResult<DashboardErrors>;
    }
  | { status: "unauthenticated" };

async function fetchFromApi(path: string): Promise<Response> {
  const cookieHeader = (await cookies()).toString();

  return fetch(`${apiOrigin}${path}`, {
    cache: "no-store",
    headers: cookieHeader ? { cookie: cookieHeader } : undefined,
  });
}

export async function getCurrentUser(): Promise<User | null> {
  try {
    const response = await fetchFromApi("/api/v1/auth/me");

    if (!response.ok) {
      return null;
    }

    const { user } = (await response.json()) as AuthResponse;
    return user ?? null;
  } catch {
    return null;
  }
}

export async function getDashboardSummary(): Promise<DashboardSummaryResult> {
  const result = await getDashboardResource<DashboardSummary>(
    "/api/v1/dashboard/summary",
  );

  if (result.status === "success") {
    return { status: "success", summary: result.data };
  }

  return result;
}

async function getDashboardResource<T>(
  path: string,
): Promise<DashboardResourceResult<T>> {
  try {
    const response = await fetchFromApi(path);

    if (response.status === 401) {
      return { status: "unauthenticated" };
    }

    if (response.status === 503) {
      return { status: "source_unavailable" };
    }

    if (!response.ok) {
      return { status: "unexpected_error" };
    }

    return {
      status: "success",
      data: (await response.json()) as T,
    };
  } catch {
    return { status: "unexpected_error" };
  }
}

export async function getDashboardOverview(): Promise<DashboardOverviewResult> {
  const [summary, trend, errors] = await Promise.all([
    getDashboardResource<DashboardSummary>("/api/v1/dashboard/summary"),
    getDashboardResource<ExecutionTrend>("/api/v1/dashboard/execution-trend"),
    getDashboardResource<DashboardErrors>("/api/v1/dashboard/errors"),
  ]);

  if (
    summary.status === "unauthenticated" ||
    trend.status === "unauthenticated" ||
    errors.status === "unauthenticated"
  ) {
    return { status: "unauthenticated" };
  }

  return { status: "success", summary, trend, errors };
}
