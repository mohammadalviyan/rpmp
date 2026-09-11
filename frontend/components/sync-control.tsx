"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";
import type { ApiError, SyncStarted, SyncStatus } from "@/lib/api/types";

const DEFAULT_POLL_INTERVAL_MS = 2000;
const DEFAULT_MAX_POLL_ATTEMPTS = 60;

type SyncState =
  | { kind: "idle" }
  | { kind: "starting" }
  | { kind: "running"; runId: string }
  | { kind: "finished" }
  | { kind: "failed" }
  | { kind: "still_running" }
  | { kind: "already_running" }
  | { kind: "forbidden" }
  | { kind: "unknown_outcome" }
  | { kind: "unavailable" };

const MESSAGES: Record<Exclude<SyncState["kind"], "idle">, string> = {
  starting: "Starting a sync of the stored RPMP copy.",
  running: "Sync is running. This updates the stored RPMP copy in the background.",
  finished: "Sync finished. The dashboard now shows the latest stored copy.",
  failed:
    "Sync did not finish. The dashboard still shows the previous stored copy.",
  still_running:
    "Sync is taking longer than usual. Check the last refresh time again in a few minutes.",
  already_running:
    "A sync is already running. Wait for it to finish, then check the last refresh time.",
  forbidden: "Your account cannot start a sync.",
  unknown_outcome:
    "Sync started, but its result could not be read. Check the last refresh time in a few minutes.",
  unavailable: "Could not start the sync. Try again in a few minutes.",
};

const ERROR_KINDS = new Set<SyncState["kind"]>([
  "failed",
  "already_running",
  "forbidden",
  "unknown_outcome",
  "unavailable",
]);

type SyncControlProps = {
  pollIntervalMs?: number;
  maxPollAttempts?: number;
};

export function SyncControl({
  pollIntervalMs = DEFAULT_POLL_INTERVAL_MS,
  maxPollAttempts = DEFAULT_MAX_POLL_ATTEMPTS,
}: SyncControlProps) {
  const router = useRouter();
  const [state, setState] = useState<SyncState>({ kind: "idle" });
  const pending = state.kind === "starting" || state.kind === "running";
  const runId = state.kind === "running" ? state.runId : null;

  useEffect(() => {
    if (runId === null) {
      return;
    }

    let cancelled = false;
    let attempts = 0;
    let timer: ReturnType<typeof setTimeout> | undefined;
    const controller = new AbortController();

    async function poll() {
      attempts += 1;

      try {
        const response = await fetch("/api/v1/sync/status", {
          credentials: "include",
          cache: "no-store",
          signal: controller.signal,
        });

        if (cancelled) {
          return;
        }

        if (!response.ok) {
          setState({ kind: "unknown_outcome" });
          return;
        }

        const status = (await response.json()) as SyncStatus;

        if (cancelled) {
          return;
        }

        if (status.run_id !== runId) {
          // Another run replaced ours as the latest, so this outcome is unknown.
          setState({ kind: "unknown_outcome" });
          router.refresh();
          return;
        }

        if (status.status === "running") {
          if (attempts >= maxPollAttempts) {
            setState({ kind: "still_running" });
            return;
          }

          timer = setTimeout(poll, pollIntervalMs);
          return;
        }

        if (status.status === "failure") {
          setState({ kind: "failed" });
          return;
        }

        setState({ kind: "finished" });
        router.refresh();
      } catch {
        if (!cancelled) {
          setState({ kind: "unknown_outcome" });
        }
      }
    }

    poll();

    return () => {
      cancelled = true;
      controller.abort();
      clearTimeout(timer);
    };
  }, [maxPollAttempts, pollIntervalMs, router, runId]);

  async function startSync() {
    if (pending) {
      return;
    }

    setState({ kind: "starting" });

    try {
      const response = await fetch("/api/v1/sync", {
        method: "POST",
        credentials: "include",
      });

      if (response.status === 202) {
        const started = (await response.json()) as SyncStarted;
        setState({ kind: "running", runId: started.run_id });
        return;
      }

      let error: Partial<ApiError> = {};
      try {
        error = (await response.json()) as Partial<ApiError>;
      } catch {
        // A non-JSON response is handled by status code alone.
      }

      if (response.status === 409 || error.code === "sync_in_progress") {
        setState({ kind: "already_running" });
      } else if (response.status === 403) {
        setState({ kind: "forbidden" });
      } else {
        setState({ kind: "unavailable" });
      }
    } catch {
      setState({ kind: "unavailable" });
    }
  }

  return (
    <div className="flex flex-col items-start gap-2 sm:items-end">
      <Button
        aria-busy={pending}
        disabled={pending}
        onClick={startSync}
        size="sm"
        type="button"
        variant="outline"
      >
        {pending ? "Syncing..." : "Sync now"}
      </Button>
      {state.kind === "idle" ? null : (
        <p
          key={state.kind}
          className={`text-xs ${
            ERROR_KINDS.has(state.kind)
              ? "text-destructive"
              : "text-muted-foreground"
          }`}
          role={ERROR_KINDS.has(state.kind) ? "alert" : "status"}
        >
          {MESSAGES[state.kind]}
        </p>
      )}
    </div>
  );
}
