import { fireEvent, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";

import {
  EmailSettingsPreview,
  GeneralSettingsPreview,
  UserManagementPreview,
} from "@/components/settings-preview-pages";

const toast = vi.hoisted(() => ({
  error: vi.fn(),
  success: vi.fn(),
}));

vi.mock("sonner", () => ({ toast }));

afterEach(() => {
  vi.clearAllMocks();
});

describe("GeneralSettingsPreview", () => {
  it("keeps changes local and labels the save as simulated", () => {
    render(<GeneralSettingsPreview />);

    expect(screen.getByText("Preview only")).toBeVisible();
    fireEvent.change(screen.getByLabelText("Workspace name"), {
      target: { value: "Local RPMP Preview" },
    });
    expect(screen.getByDisplayValue("Local RPMP Preview")).toBeVisible();

    fireEvent.click(screen.getByRole("button", { name: "Simulate save" }));
    expect(toast.success).toHaveBeenCalledWith(
      "Settings save simulated",
      expect.objectContaining({
        description: expect.stringMatching(/Nothing was saved/i),
      }),
    );
  });
});

describe("UserManagementPreview", () => {
  it("adds a preview member without creating an account or sending email", () => {
    render(<UserManagementPreview />);

    expect(screen.getByText("Preview only")).toBeVisible();
    fireEvent.change(screen.getByLabelText("Preview member email"), {
      target: { value: "preview.user@example.com" },
    });
    fireEvent.click(
      screen.getByRole("button", { name: "Simulate invite" }),
    );

    expect(screen.getByText("Sample members (5)")).toBeVisible();
    expect(screen.getByText("preview.user@example.com")).toBeVisible();
    expect(toast.success).toHaveBeenCalledWith(
      "Invitation simulated",
      expect.objectContaining({
        description: expect.stringMatching(
          /No invitation was sent and no user was created/i,
        ),
      }),
    );
  });

  it("rejects invalid preview email locally", () => {
    render(<UserManagementPreview />);

    fireEvent.change(screen.getByLabelText("Preview member email"), {
      target: { value: "invalid" },
    });
    fireEvent.click(
      screen.getByRole("button", { name: "Simulate invite" }),
    );

    expect(screen.getByText("Sample members (4)")).toBeVisible();
    expect(toast.error).toHaveBeenCalledWith(
      "Enter a valid email address.",
    );
  });
});

describe("EmailSettingsPreview", () => {
  it("simulates test and save without connecting or sending", () => {
    render(<EmailSettingsPreview />);

    expect(screen.getByText("Preview only")).toBeVisible();
    fireEvent.change(screen.getByLabelText("Host"), {
      target: { value: "local.example.test" },
    });

    fireEvent.click(
      screen.getByRole("button", { name: "Simulate test email" }),
    );
    expect(toast.success).toHaveBeenCalledWith(
      "Email test simulated",
      expect.objectContaining({
        description: expect.stringMatching(/No test email was sent/i),
      }),
    );

    fireEvent.click(screen.getByRole("button", { name: "Simulate save" }));
    expect(toast.success).toHaveBeenCalledWith(
      "Email settings save simulated",
      expect.objectContaining({
        description: expect.stringMatching(
          /local\.example\.test.*Nothing was saved/i,
        ),
      }),
    );
  });
});
