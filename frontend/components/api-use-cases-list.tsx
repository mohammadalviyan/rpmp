"use client";

import { useState } from "react";
import { ChevronRight, Search } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { Input } from "@/components/ui/input";
import { Progress } from "@/components/ui/progress";
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
  UseCaseListItem,
} from "@/lib/api/types";
import { cn } from "@/lib/utils";

const statusClasses: Record<ApiUseCaseStatus, string> = {
  active:
    "border-[var(--bri-green-300)] bg-[var(--bri-green-100)] text-[var(--bri-green-700)]",
  inactive: "border-border bg-muted text-muted-foreground",
};

function ApiStatusBadge({ status }: { status: ApiUseCaseStatus }) {
  return (
    <span
      className={cn(
        "inline-flex rounded-full border px-2.5 py-0.5 text-[11px] font-semibold capitalize",
        statusClasses[status],
      )}
    >
      {status}
    </span>
  );
}

export function ApiUseCasesList({
  useCases,
}: {
  useCases: UseCaseListItem[];
}) {
  const [query, setQuery] = useState("");
  const router = useRouter();
  const normalizedQuery = query.trim().toLowerCase();
  const rows = useCases.filter((useCase) =>
    useCase.name.toLowerCase().includes(normalizedQuery),
  );

  const navigateTo = (id: string) => router.push(`/use-cases/${id}`);

  return (
    <main className="mx-auto w-full max-w-[1600px] space-y-6 px-4 py-8 sm:px-6 lg:px-8">
      <div className="relative min-w-0">
        <Search
          aria-hidden="true"
          className="pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground"
        />
        <Input
          aria-label="Search use cases"
          className="w-full rounded-full pl-9 sm:w-72"
          onChange={(event) => setQuery(event.target.value)}
          placeholder="Search use case..."
          value={query}
        />
      </div>

      <div className="overflow-hidden rounded-2xl border bg-card shadow-[0_8px_24px_var(--bri-black-opacity-10)]">
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow className="hover:bg-transparent">
                <TableHead className="min-w-56">Use Case</TableHead>
                <TableHead className="min-w-40">Success Rate</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Failed</TableHead>
                <TableHead>Executions</TableHead>
                <TableHead>Processes</TableHead>
                <TableHead className="w-10">
                  <span className="sr-only">Open details</span>
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {rows.map((useCase) => (
                <TableRow
                  aria-label={`Open ${useCase.name}`}
                  className="cursor-pointer"
                  key={useCase.id}
                  onClick={(event) => {
                    if (!(event.target as HTMLElement).closest("a")) {
                      navigateTo(useCase.id);
                    }
                  }}
                >
                  <TableCell>
                    <Link
                      className="block min-w-0"
                      href={`/use-cases/${useCase.id}`}
                    >
                      <p className="truncate text-sm font-semibold">
                        {useCase.name}
                      </p>
                      <p className="truncate text-xs text-muted-foreground">
                        {useCase.environments.length > 0
                          ? useCase.environments.join(", ")
                          : "No environment recorded"}
                      </p>
                    </Link>
                  </TableCell>
                  <TableCell>
                    {useCase.success_rate === null ? (
                      <span className="text-xs text-muted-foreground">
                        Not available
                      </span>
                    ) : (
                      <div className="flex items-center gap-2">
                        <Progress
                          aria-label={`${useCase.name} success rate`}
                          className="w-24"
                          value={useCase.success_rate}
                        />
                        <span className="text-xs font-semibold">
                          {useCase.success_rate.toFixed(2)}%
                        </span>
                      </div>
                    )}
                  </TableCell>
                  <TableCell>
                    <ApiStatusBadge status={useCase.status} />
                  </TableCell>
                  <TableCell className="text-sm font-semibold tabular-nums">
                    {useCase.failed_count.toLocaleString()}
                  </TableCell>
                  <TableCell className="text-sm tabular-nums">
                    {useCase.execution_volume.toLocaleString()}
                  </TableCell>
                  <TableCell className="text-sm tabular-nums">
                    {useCase.process_count.toLocaleString()}
                  </TableCell>
                  <TableCell>
                    <Link
                      aria-label={`View ${useCase.name} details`}
                      href={`/use-cases/${useCase.id}`}
                    >
                      <ChevronRight
                        aria-hidden="true"
                        className="h-4 w-4 text-muted-foreground"
                      />
                    </Link>
                  </TableCell>
                </TableRow>
              ))}
              {rows.length === 0 ? (
                <TableRow>
                  <TableCell
                    className="py-10 text-center text-sm text-muted-foreground"
                    colSpan={7}
                  >
                    {useCases.length === 0
                      ? "No use cases are available in the latest snapshot."
                      : "No use case matches your search. Clear the search to see all use cases."}
                  </TableCell>
                </TableRow>
              ) : null}
            </TableBody>
          </Table>
        </div>
      </div>
    </main>
  );
}
