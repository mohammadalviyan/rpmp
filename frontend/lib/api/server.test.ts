import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("next/headers", () => ({
  cookies: async () => ({
    toString: () => "rpmp_access=http-only-cookie",
  }),
}));

const { getCurrentUser } = await import("@/lib/api/server");

const fetchMock = vi.fn();

describe("getCurrentUser", () => {
  beforeEach(() => {
    fetchMock.mockReset();
    vi.stubGlobal("fetch", fetchMock);
  });

  it("loads the current user from GET /api/v1/auth/me with the request cookie", async () => {
    fetchMock.mockResolvedValue(
      new Response(
        JSON.stringify({
          user: {
            id: "018f5f71-9cb9-7a61-97e7-d8f33f12b821",
            employee_id: "12345678",
            display_name: "Example Viewer",
            role: "viewer",
          },
        }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      ),
    );

    await expect(getCurrentUser()).resolves.toEqual({
      id: "018f5f71-9cb9-7a61-97e7-d8f33f12b821",
      employee_id: "12345678",
      display_name: "Example Viewer",
      role: "viewer",
    });
    expect(fetchMock).toHaveBeenCalledWith(
      "http://127.0.0.1:8080/api/v1/auth/me",
      expect.objectContaining({
        cache: "no-store",
        headers: { cookie: "rpmp_access=http-only-cookie" },
      }),
    );
  });

  it("returns null when me rejects the session", async () => {
    fetchMock.mockResolvedValue(new Response(null, { status: 401 }));

    await expect(getCurrentUser()).resolves.toBeNull();
  });
});
