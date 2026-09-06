import type { ReactNode } from "react";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import LoginPage from "@/app/(auth)/login/page";
import ProtectedLayout from "@/app/(protected)/layout";
import { getCurrentUser } from "@/lib/api/server";
import type { User } from "@/lib/api/types";

const redirect = vi.fn();

vi.mock("next/navigation", () => ({
  redirect: (destination: string) => redirect(destination),
  useRouter: () => ({ replace: vi.fn(), refresh: vi.fn() }),
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

  it("renders the protected Overview shell with signed-in identity", async () => {
    mockGetCurrentUser.mockResolvedValue(viewer);

    render(
      await ProtectedLayout({
        children: <p>Dashboard content</p>,
      }),
    );

    expect(
      screen.getByRole("link", { name: "RPMP Overview" }),
    ).toHaveAttribute("href", "/dashboard");
    expect(
      screen.getByRole("link", { name: "Overview" }),
    ).toHaveAttribute("aria-current", "page");
    for (const group of ["Dashboard", "Report Management", "Email", "Settings"]) {
      expect(screen.getByText(group)).toBeVisible();
    }

    const placeholders = [
      "Use Cases",
      "Generate Report",
      "Report History",
      "Email Management",
      "Email History",
      "General Settings",
      "User Management",
      "Email Configuration",
    ];
    for (const label of placeholders) {
      expect(screen.getByRole("button", { name: label })).toBeDisabled();
    }
    expect(screen.getAllByRole("link")).toHaveLength(2);
    expect(screen.getByText("Example Viewer")).toBeVisible();
    expect(screen.getByText("viewer")).toBeVisible();
    expect(screen.getByRole("button", { name: "Log out" })).toBeVisible();
  });

  it("opens and closes the mobile navigation drawer", async () => {
    mockGetCurrentUser.mockResolvedValue(viewer);
    const user = userEvent.setup();

    render(
      await ProtectedLayout({
        children: <p>Dashboard content</p>,
      }),
    );

    const trigger = screen.getByRole("button", { name: "Open menu" });
    expect(trigger).toHaveAttribute("aria-expanded", "false");
    expect(trigger).toHaveAttribute("aria-controls", "app-sidebar");

    await user.click(trigger);

    expect(
      screen.getByRole("button", { name: "Close menu", expanded: true }),
    ).toBeVisible();

    await user.keyboard("{Escape}");

    expect(
      screen.getByRole("button", { name: "Open menu" }),
    ).toHaveAttribute("aria-expanded", "false");
  });

  it("keeps the RPMP sign-in hierarchy for an unauthenticated visitor", async () => {
    mockGetCurrentUser.mockResolvedValue(null);

    render(await LoginPage());

    expect(screen.getByText("RPMP")).toBeVisible();
    expect(screen.getByRole("heading", { name: "Sign in" })).toBeVisible();
    expect(
      screen.getByText("Use your Employee ID and password to continue."),
    ).toBeVisible();
    expect(screen.getByRole("button", { name: "Sign in" })).toBeVisible();
    expect(redirect).not.toHaveBeenCalled();
  });
});
