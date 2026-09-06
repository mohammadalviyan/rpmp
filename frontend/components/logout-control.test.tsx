import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { LogoutControl } from "@/components/logout-control";

const replace = vi.fn();
const refresh = vi.fn();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ replace, refresh }),
}));

const fetchMock = vi.fn();

describe("LogoutControl", () => {
  beforeEach(() => {
    fetchMock.mockReset();
    replace.mockReset();
    refresh.mockReset();
    vi.stubGlobal("fetch", fetchMock);
  });

  it("ends the cookie session and navigates to login", async () => {
    fetchMock.mockResolvedValue(new Response(null, { status: 204 }));
    const user = userEvent.setup();
    render(<LogoutControl />);

    await user.click(screen.getByRole("button", { name: "Log out" }));

    expect(fetchMock).toHaveBeenCalledWith("/api/v1/auth/logout", {
      method: "POST",
      credentials: "include",
    });
    expect(replace).toHaveBeenCalledWith("/login");
    expect(refresh).toHaveBeenCalled();
  });

  it("keeps the user in place when logout fails", async () => {
    fetchMock.mockResolvedValue(new Response(null, { status: 503 }));
    const user = userEvent.setup();
    render(<LogoutControl />);

    await user.click(screen.getByRole("button", { name: "Log out" }));

    expect(screen.getByRole("alert")).toHaveTextContent(
      "Logout failed. Try again.",
    );
    expect(replace).not.toHaveBeenCalled();
  });
});
