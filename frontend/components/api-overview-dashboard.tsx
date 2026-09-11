import { DashboardSummary } from "@/components/dashboard-summary";
import { ErrorDistribution } from "@/components/error-distribution";
import { ExecutionTrendChart } from "@/components/execution-trend-chart";
import { SyncControl } from "@/components/sync-control";
import type {
  DashboardErrors,
  DashboardSummary as DashboardSummaryData,
  ExecutionTrend,
} from "@/lib/api/types";
import type { DashboardResourceResult } from "@/lib/api/server";

type ApiOverviewDashboardProps = {
  summary: DashboardResourceResult<DashboardSummaryData>;
  trend: DashboardResourceResult<ExecutionTrend>;
  errors: DashboardResourceResult<DashboardErrors>;
  canTriggerSync?: boolean;
};

const cardClass =
  "rounded-2xl border bg-card shadow-[0_8px_24px_var(--bri-black-opacity-10)]";

function ResourceUnavailable({
  resource,
  status,
}: {
  resource: string;
  status: "unauthenticated" | "source_unavailable" | "unexpected_error";
}) {
  if (status === "unauthenticated") {
    return null;
  }

  return (
    <p role="status" className="mt-6 text-sm text-muted-foreground">
      {status === "source_unavailable"
        ? `${resource} is unavailable because the stored RPMP data source could not be read.`
        : `${resource} could not be loaded. Try again later.`}
    </p>
  );
}

export function ApiOverviewDashboard({
  summary,
  trend,
  errors,
  canTriggerSync = false,
}: ApiOverviewDashboardProps) {
  return (
    <main
      aria-label="Overview dashboard"
      className="mx-auto w-full max-w-[1600px] space-y-8 px-4 py-8 sm:px-6 lg:px-8"
    >
      {canTriggerSync ? (
        <section aria-label="Data sync" className="flex justify-end">
          <SyncControl />
        </section>
      ) : null}

      {summary.status === "success" ? (
        <DashboardSummary summary={summary.data} />
      ) : (
        <section aria-label="Dashboard summary" className={`${cardClass} p-5`}>
          <h2 className="text-sm font-semibold">Dashboard summary</h2>
          <ResourceUnavailable resource="Dashboard summary" status={summary.status} />
        </section>
      )}

      <section className="grid gap-4 lg:grid-cols-3">
        <article className={`${cardClass} min-w-0 p-5 lg:col-span-2`}>
          <h2 className="text-sm font-semibold">Execution Trend</h2>
          <p className="text-xs text-muted-foreground">
            Successful and failed executions in the selected period
          </p>
          {trend.status === "success" ? (
            <ExecutionTrendChart trend={trend.data} />
          ) : (
            <ResourceUnavailable resource="Execution trend" status={trend.status} />
          )}
        </article>

        <article className={`${cardClass} p-5`}>
          <h2 className="text-sm font-semibold">Top Errors</h2>
          <p className="text-xs text-muted-foreground">
            Normalized error categories in the selected period
          </p>
          {errors.status === "success" ? (
            <ErrorDistribution errors={errors.data} />
          ) : (
            <ResourceUnavailable resource="Top errors" status={errors.status} />
          )}
        </article>
      </section>
    </main>
  );
}
