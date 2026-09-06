import { beforeEach, describe, expect, it, vi } from "vitest";

vi.mock("next/headers", () => ({
  cookies: async () => ({
    toString: () => "rpmp_access=http-only-cookie",
  }),
}));

const { getCurrentUser, getDashboardOverview, getDashboardSummary } =
  await import("@/lib/api/server");

const fetchMock = vi.fn();
const dashboardSummary = {
  period: {
    from: "2026-08-07T00:00:00Z",
    to: "2026-09-06T00:00:00Z",
    timezone: "UTC",
  },
  freshness: {
    status: "fresh",
    last_successful_refresh_at: "2026-09-06T00:00:00Z",
  },
  kpis: {
    total_use_cases: 12,
    active_use_cases: 9,
    execution_volume: 100,
    success_rate: 94,
    failed_executions: 6,
  },
};

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

describe("getDashboardSummary", () => {
  beforeEach(() => {
    fetchMock.mockReset();
    vi.stubGlobal("fetch", fetchMock);
  });

  it("loads the summary without caching and forwards the request cookie", async () => {
    fetchMock.mockResolvedValue(
      new Response(JSON.stringify(dashboardSummary), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      }),
    );

    await expect(getDashboardSummary()).resolves.toEqual({
      status: "success",
      summary: dashboardSummary,
    });
    expect(fetchMock).toHaveBeenCalledWith(
      "http://127.0.0.1:8080/api/v1/dashboard/summary",
      {
        cache: "no-store",
        headers: { cookie: "rpmp_access=http-only-cookie" },
      },
    );
  });

  it("classifies a 401 response as unauthenticated", async () => {
    fetchMock.mockResolvedValue(new Response(null, { status: 401 }));

    await expect(getDashboardSummary()).resolves.toEqual({
      status: "unauthenticated",
    });
  });

  it("classifies a 503 response as source unavailable", async () => {
    fetchMock.mockResolvedValue(
      new Response(
        JSON.stringify({
          code: "source_unavailable",
          message: "Source unavailable",
        }),
        { status: 503, headers: { "Content-Type": "application/json" } },
      ),
    );

    await expect(getDashboardSummary()).resolves.toEqual({
      status: "source_unavailable",
    });
  });

  it("uses a safe unexpected result for forbidden and network failures", async () => {
    fetchMock.mockResolvedValueOnce(new Response(null, { status: 403 }));

    await expect(getDashboardSummary()).resolves.toEqual({
      status: "unexpected_error",
    });

    fetchMock.mockRejectedValueOnce(new Error("connection failed"));

    await expect(getDashboardSummary()).resolves.toEqual({
      status: "unexpected_error",
    });
  });
});

describe("getDashboardOverview", () => {
  beforeEach(() => {
    fetchMock.mockReset();
    vi.stubGlobal("fetch", fetchMock);
  });

  it("loads all dashboard resources independently with no-store cookie forwarding", async () => {
    const trend = {
      period: dashboardSummary.period,
      points: [
        {
          bucket: "2026-08-01T00:00:00Z",
          label: "Aug",
          success: 94,
          failure: 6,
        },
      ],
    };
    const errors = {
      period: dashboardSummary.period,
      groups: [{ code: "faulted", label: "Faulted", count: 4 }],
    };

    fetchMock
      .mockResolvedValueOnce(
        new Response(JSON.stringify(dashboardSummary), { status: 200 }),
      )
      .mockResolvedValueOnce(
        new Response(JSON.stringify(trend), { status: 200 }),
      )
      .mockResolvedValueOnce(
        new Response(JSON.stringify(errors), { status: 200 }),
      );

    await expect(getDashboardOverview()).resolves.toEqual({
      status: "success",
      summary: { status: "success", data: dashboardSummary },
      trend: { status: "success", data: trend },
      errors: { status: "success", data: errors },
    });
    expect(fetchMock).toHaveBeenCalledTimes(3);
    for (const path of [
      "summary",
      "execution-trend",
      "errors",
    ]) {
      expect(fetchMock).toHaveBeenCalledWith(
        `http://127.0.0.1:8080/api/v1/dashboard/${path}`,
        {
          cache: "no-store",
          headers: { cookie: "rpmp_access=http-only-cookie" },
        },
      );
    }
  });

  it("redirects the whole read on 401 but keeps 503 states independent", async () => {
    fetchMock
      .mockResolvedValueOnce(new Response(null, { status: 200 }))
      .mockResolvedValueOnce(new Response(null, { status: 401 }))
      .mockResolvedValueOnce(new Response(null, { status: 503 }));

    await expect(getDashboardOverview()).resolves.toEqual({
      status: "unauthenticated",
    });

    fetchMock
      .mockResolvedValueOnce(
        new Response(JSON.stringify(dashboardSummary), { status: 200 }),
      )
      .mockResolvedValueOnce(new Response(null, { status: 503 }))
      .mockResolvedValueOnce(new Response("not-json", { status: 200 }));

    await expect(getDashboardOverview()).resolves.toEqual({
      status: "success",
      summary: { status: "success", data: dashboardSummary },
      trend: { status: "source_unavailable" },
      errors: { status: "unexpected_error" },
    });
  });
});
