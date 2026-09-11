import {
  Activity,
  ArrowLeft,
  CheckCircle2,
  ListTodo,
  XCircle,
} from "lucide-react";
import Link from "next/link";

import { buttonVariants } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type {
  ApiUseCaseStatus,
  UseCaseResponse,
} from "@/lib/api/types";
import { cn } from "@/lib/utils";

const cardClass =
  "rounded-2xl border bg-card shadow-[0_8px_24px_var(--bri-black-opacity-10)]";

const utcDateTimeFormatter = new Intl.DateTimeFormat("en-US", {
  dateStyle: "medium",
  timeStyle: "short",
  timeZone: "UTC",
});

function formatUtcDateTime(value: string) {
  return utcDateTimeFormatter.format(new Date(value));
}

function ApiStatusBadge({ status }: { status: ApiUseCaseStatus }) {
  return (
    <span
      className={cn(
        "inline-flex rounded-full border px-2.5 py-0.5 text-[11px] font-semibold capitalize",
        status === "active"
          ? "border-[var(--bri-green-300)] bg-[var(--bri-green-100)] text-[var(--bri-green-700)]"
          : "border-border bg-muted text-muted-foreground",
      )}
    >
      {status}
    </span>
  );
}

export function ApiUseCaseDetail({ data }: { data: UseCaseResponse }) {
  const useCase = data.use_case;
  const stats = [
    {
      label: "Execution Volume",
      value: useCase.execution_volume.toLocaleString(),
      icon: ListTodo,
    },
    {
      label: "Successful",
      value: useCase.successful_count.toLocaleString(),
      icon: CheckCircle2,
    },
    {
      label: "Failed",
      value: useCase.failed_count.toLocaleString(),
      icon: XCircle,
    },
    {
      label: "Success Rate",
      value:
        useCase.success_rate === null
          ? "Not available"
          : `${useCase.success_rate.toFixed(2)}%`,
      icon: Activity,
    },
  ];

  return (
    <main className="mx-auto w-full max-w-[1600px] space-y-6 px-4 py-8 sm:px-6 lg:px-8">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <Link
          className={cn(
            buttonVariants({ variant: "outline", size: "sm" }),
            "rounded-full",
          )}
          href="/use-cases"
        >
          <ArrowLeft aria-hidden="true" className="mr-1 h-4 w-4" />
          Back to Use Cases
        </Link>
        <ApiStatusBadge status={useCase.status} />
      </div>

      <section aria-labelledby="use-case-name">
        <h1 id="use-case-name" className="text-2xl font-bold tracking-tight">
          {useCase.name}
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Source key: {useCase.source_key}
        </p>
        <p className="mt-1 text-sm text-muted-foreground">
          Environments:{" "}
          {useCase.environments.length > 0
            ? useCase.environments.join(", ")
            : "None recorded"}
        </p>
        <p className="mt-2 text-xs text-muted-foreground">
          This is a stored RPMP copy. Last successful refresh:{" "}
          {formatUtcDateTime(data.freshness.last_successful_refresh_at)} UTC.
          Status: {data.freshness.status}. Use case updated:{" "}
          {formatUtcDateTime(useCase.updated_at)} UTC.
        </p>
      </section>

      <section
        aria-label="Use case key performance indicators"
        className="grid grid-cols-2 gap-4 lg:grid-cols-4"
      >
        {stats.map((stat) => (
          <article className={`${cardClass} p-5`} key={stat.label}>
            <div className="flex items-center justify-between">
              <h2 className="text-xs text-muted-foreground">{stat.label}</h2>
              <stat.icon aria-hidden="true" className="h-4 w-4 text-primary" />
            </div>
            <p className="mt-2 text-2xl font-bold tracking-tight tabular-nums">
              {stat.value}
            </p>
          </article>
        ))}
      </section>

      <section
        className={`${cardClass} overflow-hidden`}
        aria-labelledby="processes"
      >
        <div className="border-b border-border px-6 py-4">
          <h2 className="text-sm font-bold" id="processes">
            Processes ({useCase.process_count})
          </h2>
          <p className="text-xs text-muted-foreground">
            Counts from the latest stored snapshot
          </p>
        </div>
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow className="hover:bg-transparent">
                <TableHead className="min-w-52">Process</TableHead>
                <TableHead className="min-w-40">Package</TableHead>
                <TableHead>Environment</TableHead>
                <TableHead>Successful</TableHead>
                <TableHead>Error</TableHead>
                <TableHead>Stopped</TableHead>
                <TableHead>Failed</TableHead>
                <TableHead>Executing</TableHead>
                <TableHead>Pending</TableHead>
                <TableHead>Suspended</TableHead>
                <TableHead>Resumed</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {useCase.processes.map((process) => (
                <TableRow key={process.source_process_key}>
                  <TableCell>
                    <p className="text-sm font-semibold">
                      {process.process_name}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {process.source_process_key}
                    </p>
                  </TableCell>
                  <TableCell className="text-xs">
                    {process.package_name || "Not recorded"}
                  </TableCell>
                  <TableCell className="text-xs">
                    {process.environment_name || "Not recorded"}
                  </TableCell>
                  {[
                    process.successful_count,
                    process.error_count,
                    process.stopped_count,
                    process.failed_count,
                    process.executing_count,
                    process.pending_count,
                    process.suspended_count,
                    process.resumed_count,
                  ].map((count, index) => (
                    <TableCell
                      className="text-sm tabular-nums"
                      key={`${process.source_process_key}-${index}`}
                    >
                      {count.toLocaleString()}
                    </TableCell>
                  ))}
                </TableRow>
              ))}
              {useCase.processes.length === 0 ? (
                <TableRow>
                  <TableCell
                    className="py-10 text-center text-sm text-muted-foreground"
                    colSpan={11}
                  >
                    No processes are available in the latest snapshot.
                  </TableCell>
                </TableRow>
              ) : null}
            </TableBody>
          </Table>
        </div>
      </section>

      <section className={`${cardClass} p-6`} aria-labelledby="weekly-performance">
        <h2 className="text-sm font-bold" id="weekly-performance">
          Weekly Performance
        </h2>
        <p className="mt-2 text-sm text-muted-foreground">
          This snapshot has no execution history.
        </p>
      </section>

      <section className={`${cardClass} p-6`} aria-labelledby="issues">
        <h2 className="text-sm font-bold" id="issues">
          Issues
        </h2>
        <p className="mt-2 text-sm text-muted-foreground">
          This snapshot has no execution history.
        </p>
      </section>
    </main>
  );
}
