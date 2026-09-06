import { redirect } from "next/navigation";

import { DashboardSummary } from "@/components/dashboard-summary";
import { ErrorDistribution } from "@/components/error-distribution";
import { ExecutionTrendChart } from "@/components/execution-trend-chart";
import {
  getDashboardOverview,
  type DashboardResourceResult,
} from "@/lib/api/server";

export default async function DashboardPage() {
  const result = await getDashboardOverview();

  if (result.status === "unauthenticated") {
    redirect("/login");
  }

  return (
    <main
      aria-label="Overview dashboard"
      className="mx-auto w-full max-w-[1600px] px-5 py-6 sm:px-8 sm:py-8"
    >
      {result.summary.status === "success" ? (
        <DashboardSummary summary={result.summary.data} />
      ) : (
        <UnavailablePanel
          result={result.summary}
          title="Dashboard data unavailable"
          sourceCopy="Current dashboard data is temporarily unavailable. KPI values, period, and refresh time are not shown."
        />
      )}

      <div className="mt-6 grid gap-6 xl:grid-cols-[minmax(0,2fr)_minmax(19rem,1fr)]">
        <section
          aria-labelledby="execution-trend-title"
          className="min-w-0 rounded-2xl border bg-card p-5 shadow-[0_8px_24px_var(--bri-black-opacity-10)] sm:p-6"
        >
          <h2 id="execution-trend-title" className="text-base font-semibold">
            Execution Trend
          </h2>
          <p className="mt-1 text-sm text-muted-foreground">
            Successful and failed executions over time
          </p>
          {result.trend.status === "success" ? (
            <ExecutionTrendChart trend={result.trend.data} />
          ) : (
            <InlineUnavailable result={result.trend} resource="trend data" />
          )}
        </section>

        <section
          aria-labelledby="error-distribution-title"
          className="rounded-2xl border bg-card p-5 shadow-[0_8px_24px_var(--bri-black-opacity-10)] sm:p-6"
        >
          <h2 id="error-distribution-title" className="text-base font-semibold">
            Error Distribution
          </h2>
          <p className="mt-1 text-sm text-muted-foreground">
            Errors grouped by category
          </p>
          {result.errors.status === "success" ? (
            <ErrorDistribution errors={result.errors.data} />
          ) : (
            <InlineUnavailable
              result={result.errors}
              resource="error distribution"
            />
          )}
        </section>
      </div>
    </main>
  );
}

function UnavailablePanel({
  result,
  title,
  sourceCopy,
}: {
  result: Exclude<DashboardResourceResult<unknown>, { status: "success" }>;
  title: string;
  sourceCopy: string;
}) {
  return (
    <section
      aria-labelledby="dashboard-unavailable-title"
      className="rounded-2xl border bg-card p-6"
    >
      <h2 id="dashboard-unavailable-title" className="text-lg font-semibold">
        {title}
      </h2>
      <p className="mt-2 text-sm text-muted-foreground">
        {result.status === "source_unavailable"
          ? sourceCopy
          : "We could not load current dashboard data. Try again later."}
      </p>
    </section>
  );
}

function InlineUnavailable({
  result,
  resource,
}: {
  result: Exclude<DashboardResourceResult<unknown>, { status: "success" }>;
  resource: string;
}) {
  return (
    <p role="status" className="mt-8 text-sm text-muted-foreground">
      {result.status === "source_unavailable"
        ? `Current ${resource} is temporarily unavailable.`
        : `We could not load ${resource}. Try again later.`}
    </p>
  );
}
