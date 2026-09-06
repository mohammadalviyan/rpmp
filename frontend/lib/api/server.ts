import "server-only";

import { cookies } from "next/headers";

import type {
  AuthResponse,
  DashboardSummary,
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
  try {
    const response = await fetchFromApi("/api/v1/dashboard/summary");

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
      summary: (await response.json()) as DashboardSummary,
    };
  } catch {
    return { status: "unexpected_error" };
  }
}
