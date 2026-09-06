import { fireEvent, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { LoginForm } from "@/components/login-form";

const replace = vi.fn();
const refresh = vi.fn();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ replace, refresh }),
}));

const fetchMock = vi.fn();

function response(body: object, status: number) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}

async function fillAndSubmit() {
  const user = userEvent.setup();
  await user.type(screen.getByLabelText("Employee ID"), "12345678");
  await user.type(screen.getByLabelText("Password"), "user supplied secret");
  await user.click(screen.getByRole("button", { name: "Sign in" }));
}

describe("LoginForm", () => {
  beforeEach(() => {
    fetchMock.mockReset();
    replace.mockReset();
    refresh.mockReset();
    vi.stubGlobal("fetch", fetchMock);
  });

  it("has accessible Employee ID and password labels", () => {
    render(<LoginForm />);

    expect(screen.getByLabelText("Employee ID")).toHaveAttribute("type", "text");
    expect(screen.getByLabelText("Password")).toHaveAttribute(
      "type",
      "password",
    );
  });

  it("disables the form while a login is pending", () => {
    fetchMock.mockReturnValue(new Promise(() => undefined));
    render(<LoginForm />);

    fireEvent.submit(screen.getByRole("button", { name: "Sign in" }));

    expect(
      screen.getByRole("button", { name: "Signing in..." }),
    ).toBeDisabled();
    expect(screen.getByLabelText("Employee ID")).toBeDisabled();
    expect(screen.getByLabelText("Password")).toBeDisabled();
    expect(fetchMock).toHaveBeenCalledTimes(1);
  });

  it("shows safe invalid-credential copy", async () => {
    fetchMock.mockResolvedValue(
      response(
        {
          code: "invalid_credentials",
          message: "Employee ID or password is incorrect.",
        },
        401,
      ),
    );
    render(<LoginForm />);

    await fillAndSubmit();

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Employee ID or password is incorrect.",
    );
  });

  it("shows unavailable copy after a network failure", async () => {
    fetchMock.mockRejectedValue(new TypeError("network failure"));
    render(<LoginForm />);

    await fillAndSubmit();

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "RPMP is unavailable right now.",
    );
  });

  it("shows unexpected-error copy for other API failures", async () => {
    fetchMock.mockResolvedValue(
      response(
        {
          code: "account_inactive",
          message: "Account is inactive.",
        },
        403,
      ),
    );
    render(<LoginForm />);

    await fillAndSubmit();

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "We could not sign you in.",
    );
  });

  it("posts credentials with cookies and redirects without storing a token", async () => {
    const storageSet = vi.spyOn(Storage.prototype, "setItem");
    fetchMock.mockResolvedValue(
      response(
        {
          user: {
            id: "018f5f71-9cb9-7a61-97e7-d8f33f12b821",
            employee_id: "12345678",
            display_name: "Example Viewer",
            role: "viewer",
          },
        },
        200,
      ),
    );
    render(<LoginForm />);

    await fillAndSubmit();

    expect(fetchMock).toHaveBeenCalledWith("/api/v1/auth/login", {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        employee_id: "12345678",
        password: "user supplied secret",
      }),
    });
    expect(replace).toHaveBeenCalledWith("/dashboard");
    expect(refresh).toHaveBeenCalled();
    expect(storageSet).not.toHaveBeenCalled();
    expect(window.localStorage).toHaveLength(0);
    expect(window.sessionStorage).toHaveLength(0);
  });
});
