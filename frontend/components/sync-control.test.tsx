import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { SyncControl } from "@/components/sync-control";

const refresh = vi.fn();

vi.mock("next/navigation", () => ({
  useRouter: () => ({ refresh }),
}));

const fetchMock = vi.fn();

function jsonResponse(body: unknown, status: number) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}

function syncStatus(
  status: "running" | "success" | "failure" | "never",
  runId: string | null,
) {
  return jsonResponse(
    {
      run_id: runId,
      status,
      started_at: runId === null ? null : "2026-09-11T00:00:00Z",
      finished_at: status === "running" || status === "never" ? null : "2026-09-11T00:01:00Z",
      rows_written: status === "success" ? 100 : 0,
    },
    200,
  );
}

describe("SyncControl", () => {
  beforeEach(() => {
    fetchMock.mockReset();
    refresh.mockReset();
    vi.stubGlobal("fetch", fetchMock);
  });

  it("starts a sync, polls until the run finishes, then refreshes the dashboard", async () => {
    fetchMock
      .mockResolvedValueOnce(
        jsonResponse({ run_id: "run-1", status: "running" }, 202),
      )
      .mockResolvedValueOnce(syncStatus("running", "run-1"))
      .mockResolvedValueOnce(syncStatus("success", "run-1"));
    const user = userEvent.setup();

    render(<SyncControl maxPollAttempts={5} pollIntervalMs={0} />);

    await user.click(screen.getByRole("button", { name: "Sync now" }));

    expect(fetchMock).toHaveBeenNthCalledWith(1, "/api/v1/sync", {
      method: "POST",
      credentials: "include",
    });

    expect(
      await screen.findByText(
        "Sync finished. The dashboard now shows the latest stored copy.",
      ),
    ).toBeVisible();
    expect(refresh).toHaveBeenCalledOnce();
    expect(fetchMock).toHaveBeenCalledTimes(3);
    expect(fetchMock.mock.calls[1][0]).toBe("/api/v1/sync/status");
    expect(fetchMock.mock.calls[1][1]).toMatchObject({
      credentials: "include",
      cache: "no-store",
    });
    expect(screen.getByRole("button", { name: "Sync now" })).toBeEnabled();
  });

  it("disables the control while the request is in flight", async () => {
    fetchMock.mockReturnValueOnce(new Promise(() => {}));
    const user = userEvent.setup();

    render(<SyncControl pollIntervalMs={0} />);

    await user.click(screen.getByRole("button", { name: "Sync now" }));

    const button = await screen.findByRole("button", { name: "Syncing..." });
    expect(button).toBeDisabled();
    expect(button).toHaveAttribute("aria-busy", "true");
    expect(screen.getByRole("status")).toHaveTextContent(
      "Starting a sync of the stored RPMP copy.",
    );
  });

  it("stops polling after the attempt budget and does not claim the data changed", async () => {
    fetchMock
      .mockResolvedValueOnce(
        jsonResponse({ run_id: "run-1", status: "running" }, 202),
      )
      .mockImplementation(async () => syncStatus("running", "run-1"));
    const user = userEvent.setup();

    render(<SyncControl maxPollAttempts={2} pollIntervalMs={0} />);

    await user.click(screen.getByRole("button", { name: "Sync now" }));

    expect(
      await screen.findByText(
        "Sync is taking longer than usual. Check the last refresh time again in a few minutes.",
      ),
    ).toBeVisible();
    expect(refresh).not.toHaveBeenCalled();
    expect(fetchMock).toHaveBeenCalledTimes(3);
  });

  it("explains that a sync is already running on 409", async () => {
    fetchMock.mockResolvedValueOnce(
      jsonResponse({ code: "sync_in_progress", message: "sync in progress" }, 409),
    );
    const user = userEvent.setup();

    render(<SyncControl pollIntervalMs={0} />);

    await user.click(screen.getByRole("button", { name: "Sync now" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "A sync is already running. Wait for it to finish, then check the last refresh time.",
    );
    expect(fetchMock).toHaveBeenCalledOnce();
    expect(refresh).not.toHaveBeenCalled();
  });

  it("reports a failed run without refreshing the dashboard", async () => {
    fetchMock
      .mockResolvedValueOnce(
        jsonResponse({ run_id: "run-2", status: "running" }, 202),
      )
      .mockResolvedValueOnce(syncStatus("failure", "run-2"));
    const user = userEvent.setup();

    render(<SyncControl pollIntervalMs={0} />);

    await user.click(screen.getByRole("button", { name: "Sync now" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Sync did not finish. The dashboard still shows the previous stored copy.",
    );
    expect(refresh).not.toHaveBeenCalled();
  });

  it("does not claim success when another run replaced ours", async () => {
    fetchMock
      .mockResolvedValueOnce(
        jsonResponse({ run_id: "run-4", status: "running" }, 202),
      )
      .mockResolvedValueOnce(syncStatus("running", "run-5"));
    const user = userEvent.setup();

    render(<SyncControl pollIntervalMs={0} />);

    await user.click(screen.getByRole("button", { name: "Sync now" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Sync started, but its result could not be read. Check the last refresh time in a few minutes.",
    );
    expect(refresh).toHaveBeenCalledOnce();
  });

  it("keeps the dashboard usable when the request fails", async () => {
    fetchMock.mockRejectedValueOnce(new Error("network down"));
    const user = userEvent.setup();

    render(<SyncControl pollIntervalMs={0} />);

    await user.click(screen.getByRole("button", { name: "Sync now" }));

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Could not start the sync. Try again in a few minutes.",
    );
    expect(screen.getByRole("button", { name: "Sync now" })).toBeEnabled();
  });

  it("stops polling when the component unmounts", async () => {
    fetchMock
      .mockResolvedValueOnce(
        jsonResponse({ run_id: "run-3", status: "running" }, 202),
      )
      .mockImplementation(async () => syncStatus("running", "run-3"));
    const user = userEvent.setup();

    const view = render(<SyncControl maxPollAttempts={50} pollIntervalMs={5} />);

    await user.click(screen.getByRole("button", { name: "Sync now" }));
    await screen.findByRole("button", { name: "Syncing..." });
    view.unmount();

    const callsAtUnmount = fetchMock.mock.calls.length;
    await new Promise((resolve) => setTimeout(resolve, 40));

    expect(fetchMock.mock.calls.length).toBe(callsAtUnmount);
  });
});
