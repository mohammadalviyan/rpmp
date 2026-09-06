import type { ReactNode } from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import LoginPage from "@/app/(auth)/login/page";
import ProtectedLayout from "@/app/(protected)/layout";
import { getCurrentUser } from "@/lib/api/server";
import type { User } from "@/lib/api/types";

const redirect = vi.fn();

vi.mock("next/navigation", () => ({
  redirect: (destination: string) => redirect(destination),
}));

vi.mock("@/lib/api/server", () => ({
  getCurrentUser: vi.fn(),
}));

const mockGetCurrentUser = vi.mocked(getCurrentUser);
const viewer: User = {
  id: "018f5f71-9cb9-7a61-97e7-d8f33f12b821",
  employee_id: "12345678",
  display_name: "Example Viewer",
  role: "viewer",
};

function throwRedirect() {
  redirect.mockImplementation(() => {
    throw new Error("NEXT_REDIRECT");
  });
}

describe("auth route gating", () => {
  beforeEach(() => {
    redirect.mockReset();
    mockGetCurrentUser.mockReset();
    throwRedirect();
  });

  it("redirects an unauthenticated protected route to login", async () => {
    mockGetCurrentUser.mockResolvedValue(null);

    await expect(
      ProtectedLayout({ children: null as ReactNode }),
    ).rejects.toThrow("NEXT_REDIRECT");
    expect(redirect).toHaveBeenCalledWith("/login");
  });

  it("redirects an authenticated user away from login to Dashboard", async () => {
    mockGetCurrentUser.mockResolvedValue(viewer);

    await expect(LoginPage()).rejects.toThrow("NEXT_REDIRECT");
    expect(redirect).toHaveBeenCalledWith("/dashboard");
  });
});
