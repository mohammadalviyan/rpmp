import type { DashboardSummary as DashboardSummaryData } from "@/lib/api/types";

const utcDateTimeFormatter = new Intl.DateTimeFormat("en-US", {
  dateStyle: "medium",
  timeStyle: "short",
  timeZone: "UTC",
});

function formatUtcDateTime(value: string): string {
  return utcDateTimeFormatter.format(new Date(value));
}

export function DashboardSummary({
  summary,
}: {
  summary: DashboardSummaryData;
}) {
  const kpis = [
    {
      id: "total-use-cases",
      label: "Total Use Cases",
      value: summary.kpis.total_use_cases.toLocaleString("en-US"),
    },
    {
      id: "active-use-cases",
      label: "Active Use Cases",
      value: summary.kpis.active_use_cases.toLocaleString("en-US"),
    },
    {
      id: "execution-volume",
      label: "Execution Volume",
      value: summary.kpis.execution_volume.toLocaleString("en-US"),
    },
    {
      id: "success-rate",
      label: "Success Rate",
      value:
        summary.kpis.success_rate === null
          ? "N/A"
          : `${summary.kpis.success_rate}%`,
    },
    {
      id: "failed-executions",
      label: "Failed Executions",
      value: summary.kpis.failed_executions.toLocaleString("en-US"),
    },
  ];

  return (
    <>
      <div className="mt-6 flex flex-col gap-1 text-sm text-muted-foreground">
        <p>
          Period: {formatUtcDateTime(summary.period.from)} to{" "}
          {formatUtcDateTime(summary.period.to)} ({summary.period.timezone})
        </p>
        <p>
          Last refreshed:{" "}
          {formatUtcDateTime(summary.freshness.last_successful_refresh_at)} UTC.
          Freshness: {summary.freshness.status}.
        </p>
      </div>

      <section
        aria-label="Dashboard key performance indicators"
        className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-5"
      >
        {kpis.map((kpi) => (
          <article
            key={kpi.id}
            aria-labelledby={`${kpi.id}-label`}
            className="rounded-xl border bg-card p-5 text-card-foreground shadow-sm"
          >
            <h2
              id={`${kpi.id}-label`}
              className="text-sm font-medium text-muted-foreground"
            >
              {kpi.label}
            </h2>
            <p className="mt-3 text-3xl font-semibold tabular-nums">
              {kpi.value}
            </p>
          </article>
        ))}
      </section>

      {summary.kpis.execution_volume === 0 ? (
        <p
          role="status"
          className="mt-6 rounded-lg border bg-background px-4 py-3 text-sm"
        >
          There were no executions in this period.
        </p>
      ) : null}
    </>
  );
}
