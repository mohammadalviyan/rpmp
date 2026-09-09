import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { ProtectedHeader } from "@/components/protected-header";
import type { User } from "@/lib/api/types";

let pathname = "/dashboard";

vi.mock("next/navigation", () => ({
  usePathname: () => pathname,
}));

vi.mock("@/components/dashboard-header", () => ({
  DashboardHeader: ({
    title,
    subtitle,
  }: {
    title: string;
    subtitle?: string;
  }) => (
    <header>
      <h1>{title}</h1>
      <p>{subtitle}</p>
    </header>
  ),
}));

const user: User = {
  id: "user-1",
  employee_id: "12345678",
  display_name: "Example Viewer",
  role: "viewer",
};

describe("ProtectedHeader", () => {
  beforeEach(() => {
    pathname = "/dashboard";
  });

  it.each([
    ["/email-management", "Email Management"],
    ["/email-history", "Email History"],
    ["/settings/general", "General Settings"],
    ["/settings/users", "User Management"],
    ["/settings/email", "Email Configuration"],
  ])("shows the matching route title for %s", (route, title) => {
    pathname = route;

    render(<ProtectedHeader user={user} />);

    expect(screen.getByRole("heading", { name: title })).toBeVisible();
  });
});
