import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import EmailSettingsPage, {
  metadata as emailMetadata,
} from "@/app/(protected)/settings/email/page";
import GeneralSettingsPage, {
  metadata as generalMetadata,
} from "@/app/(protected)/settings/general/page";
import UserManagementPage, {
  metadata as usersMetadata,
} from "@/app/(protected)/settings/users/page";

vi.mock("@/components/settings-preview-pages", () => ({
  EmailSettingsPreview: () => <p>Email settings preview route</p>,
  GeneralSettingsPreview: () => <p>General settings preview route</p>,
  UserManagementPreview: () => <p>User management preview route</p>,
}));

describe("Settings route pages", () => {
  it("renders the general settings preview route", () => {
    render(<GeneralSettingsPage />);

    expect(screen.getByText("General settings preview route")).toBeVisible();
    expect(generalMetadata.title).toBe("General Settings Preview | RPMP");
  });

  it("renders the user management preview route", () => {
    render(<UserManagementPage />);

    expect(screen.getByText("User management preview route")).toBeVisible();
    expect(usersMetadata.title).toBe("User Management Preview | RPMP");
  });

  it("renders the email configuration preview route", () => {
    render(<EmailSettingsPage />);

    expect(screen.getByText("Email settings preview route")).toBeVisible();
    expect(emailMetadata.title).toBe("Email Configuration Preview | RPMP");
  });
});
