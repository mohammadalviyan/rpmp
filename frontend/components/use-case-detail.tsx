import {
  Activity,
  ArrowLeft,
  CheckCircle2,
  ListTodo,
  XCircle,
} from "lucide-react";
import Link from "next/link";

import {
  ErrorTypeBadge,
  IssueStatusBadge,
  StatusBadge,
  TypeBadge,
} from "@/components/rpa-badges";
import { UseCaseWeeklyChart } from "@/components/use-case-weekly-chart";
import { buttonVariants } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { UseCase } from "@/lib/data/types";
import { cn } from "@/lib/utils";

const cardClass =
  "rounded-2xl border bg-card shadow-[0_8px_24px_var(--bri-black-opacity-10)]";

export function UseCaseDetail({ useCase }: { useCase: UseCase }) {
  const stats = [
    {
      label: "Total Process",
      value: useCase.totalProcess.toLocaleString(),
      icon: ListTodo,
    },
    {
      label: "Success",
      value: useCase.totalSuccess.toLocaleString(),
      icon: CheckCircle2,
    },
    {
      label: "Failed",
      value: useCase.totalFailed.toLocaleString(),
      icon: XCircle,
    },
    {
      label: "Success Rate",
      value: `${useCase.successRate}%`,
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
        <div className="flex items-center gap-2">
          <TypeBadge type={useCase.automationType} />
          <StatusBadge status={useCase.status} />
        </div>
      </div>

      <section aria-labelledby="use-case-name">
        <h1 id="use-case-name" className="text-2xl font-bold tracking-tight">
          {useCase.name}
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          {useCase.description}
        </p>
        <p className="mt-2 text-xs font-medium text-muted-foreground">
          Business unit: {useCase.owner}
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

      <section className={`${cardClass} p-6`} aria-labelledby="weekly-performance">
        <div className="mb-4">
          <h2 className="text-sm font-bold" id="weekly-performance">
            Weekly Performance
          </h2>
          <p className="text-xs text-muted-foreground">
            Successful vs failed transactions per day
          </p>
        </div>
        <UseCaseWeeklyChart trend={useCase.trend} />
      </section>

      <section className={`${cardClass} overflow-hidden`} aria-labelledby="issues">
        <div className="border-b border-border px-6 py-4">
          <h2 className="text-sm font-bold" id="issues">
            Issues ({useCase.issues.length})
          </h2>
          <p className="text-xs text-muted-foreground">
            Errors detected for this use case
          </p>
        </div>
        <Table>
          <TableHeader>
            <TableRow className="hover:bg-transparent">
              <TableHead className="min-w-56">Issue</TableHead>
              <TableHead>Occurred At</TableHead>
              <TableHead>Error Type</TableHead>
              <TableHead>Status</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {useCase.issues.map((issue) => (
              <TableRow key={issue.id}>
                <TableCell>
                  <p className="text-sm font-semibold">{issue.name}</p>
                  <p className="max-w-xl whitespace-normal text-xs text-muted-foreground">
                    {issue.description}
                  </p>
                </TableCell>
                <TableCell className="text-xs whitespace-nowrap text-muted-foreground">
                  {issue.occurredAt}
                </TableCell>
                <TableCell>
                  <ErrorTypeBadge type={issue.errorType} />
                </TableCell>
                <TableCell>
                  <IssueStatusBadge status={issue.status} />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </section>
    </main>
  );
}

export function UseCaseNotFound() {
  return (
    <main className="mx-auto grid min-h-[50vh] w-full max-w-2xl place-items-center px-4 py-12 text-center">
      <div>
        <p className="text-sm font-semibold text-primary">Use case not found</p>
        <h1 className="mt-2 text-2xl font-bold tracking-tight">
          This use case does not exist
        </h1>
        <p className="mt-2 text-sm text-muted-foreground">
          The ID may be invalid, or the use case may no longer be available.
        </p>
        <Link
          className={cn(
            buttonVariants({ variant: "outline", size: "sm" }),
            "mt-5 rounded-full",
          )}
          href="/use-cases"
        >
          <ArrowLeft aria-hidden="true" className="mr-1 h-4 w-4" />
          Back to Use Cases
        </Link>
      </div>
    </main>
  );
}
