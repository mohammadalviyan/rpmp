"use client";

import { useState } from "react";
import { ChevronRight, Search } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { StatusBadge, TypeBadge } from "@/components/rpa-badges";
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
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import type { UseCase } from "@/lib/data/types";

type Filter = "all" | "attended" | "unattended";

export function UseCasesList({ useCases }: { useCases: UseCase[] }) {
  const [query, setQuery] = useState("");
  const [filter, setFilter] = useState<Filter>("all");
  const router = useRouter();

  const rows = useCases.filter((useCase) => {
    const matchQuery = `${useCase.name}${useCase.owner}`
      .toLowerCase()
      .includes(query.toLowerCase());
    const matchType =
      filter === "all" ||
      useCase.automationType.toLowerCase() === filter;

    return matchQuery && matchType;
  });

  const navigateTo = (id: string) => router.push(`/use-cases/${id}`);

  return (
    <main className="mx-auto w-full max-w-[1600px] space-y-6 px-4 py-8 sm:px-6 lg:px-8">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
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
        <Tabs
          className="max-w-full overflow-x-auto"
          onValueChange={(value) => setFilter(value as Filter)}
          value={filter}
        >
          <TabsList className="rounded-full">
            <TabsTrigger className="rounded-full text-xs" value="all">
              All
            </TabsTrigger>
            <TabsTrigger className="rounded-full text-xs" value="unattended">
              Unattended
            </TabsTrigger>
            <TabsTrigger className="rounded-full text-xs" value="attended">
              Attended
            </TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      <div className="overflow-hidden rounded-2xl border bg-card shadow-[0_8px_24px_var(--bri-black-opacity-10)]">
        <div className="overflow-x-auto">
          <Table>
          <TableHeader>
            <TableRow className="hover:bg-transparent">
              <TableHead className="min-w-48">Use Case</TableHead>
              <TableHead className="min-w-40">Success Rate</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Issues</TableHead>
              <TableHead>Automation Type</TableHead>
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
                  <Link className="block min-w-0" href={`/use-cases/${useCase.id}`}>
                    <p className="truncate text-sm font-semibold">{useCase.name}</p>
                    <p className="truncate text-xs text-muted-foreground">
                      {useCase.owner}
                    </p>
                  </Link>
                </TableCell>
                <TableCell>
                  <div className="flex items-center gap-2">
                    <Progress
                      aria-label={`${useCase.name} success rate`}
                      className="w-24"
                      value={useCase.successRate}
                    />
                    <span className="text-xs font-semibold">
                      {useCase.successRate}%
                    </span>
                  </div>
                </TableCell>
                <TableCell>
                  <StatusBadge status={useCase.status} />
                </TableCell>
                <TableCell className="text-sm font-semibold">
                  {useCase.issues.length}
                </TableCell>
                <TableCell>
                  <TypeBadge type={useCase.automationType} />
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
                  colSpan={6}
                >
                  No use case matches your search.
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
