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
  usePathname: () => "/dashboard",
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

    const links = {
      "Use Cases": "/use-cases",
      "Generate Report": "/generate-report",
      "Report History": "/report-history",
      "Email Management": "/email-management",
      "Email History": "/email-history",
      "General Settings": "/settings/general",
      "User Management": "/settings/users",
      "Email Configuration": "/settings/email",
    };
    for (const [label, href] of Object.entries(links)) {
      expect(screen.getByRole("link", { name: label })).toHaveAttribute(
        "href",
        href,
      );
    }
    expect(screen.getAllByRole("link")).toHaveLength(10);
    expect(screen.getByText("Example Viewer")).toBeVisible();
    expect(screen.getByText("viewer")).toBeVisible();
    expect(screen.getByRole("status")).toHaveTextContent("Demo data");
    expect(screen.getByRole("button", { name: "Log out" })).toBeVisible();
  });

  it("collapses and expands the sidebar", async () => {
    mockGetCurrentUser.mockResolvedValue(viewer);
    const user = userEvent.setup();

    render(
      await ProtectedLayout({
        children: <p>Dashboard content</p>,
      }),
    );

    const trigger = screen.getByRole("button", { name: "Toggle sidebar" });
    expect(trigger).toHaveAttribute("aria-expanded", "true");

    await user.click(trigger);
    expect(trigger).toHaveAttribute("aria-expanded", "false");
    expect(document.cookie).toContain("sidebar_state=false");

    await user.click(trigger);
    expect(trigger).toHaveAttribute("aria-expanded", "true");
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
