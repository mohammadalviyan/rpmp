import "server-only";

import { cookies } from "next/headers";

import type { AuthResponse, User } from "@/lib/api/types";

const apiOrigin = (
  process.env.RPMP_API_ORIGIN ?? "http://127.0.0.1:8080"
).replace(/\/$/, "");

export async function getCurrentUser(): Promise<User | null> {
  try {
    const cookieHeader = (await cookies()).toString();
    const response = await fetch(`${apiOrigin}/api/v1/auth/me`, {
      cache: "no-store",
      headers: cookieHeader ? { cookie: cookieHeader } : undefined,
    });

    if (!response.ok) {
      return null;
    }

    const { user } = (await response.json()) as AuthResponse;
    return user ?? null;
  } catch {
    return null;
  }
}
