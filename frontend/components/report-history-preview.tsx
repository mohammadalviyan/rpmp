"use client";

import { Download, Search } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

import { ReportStatusBadge } from "@/components/rpa-badges";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { ReportRecord } from "@/lib/data/types";

export function ReportHistoryPreview({
  reportHistory,
}: {
  reportHistory: ReportRecord[];
}) {
  const [query, setQuery] = useState("");
  const rows = reportHistory.filter((report) =>
    `${report.id}${report.useCase}${report.type}`
      .toLowerCase()
      .includes(query.toLowerCase()),
  );

  return (
    <main className="flex-1 space-y-6 px-4 py-8 sm:px-6 lg:px-8">
      <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
        <div className="flex items-center gap-2">
          <span className="rounded-full border border-primary/30 bg-primary/10 px-2.5 py-1 text-xs font-semibold text-primary">
            Demo data
          </span>
          <p className="text-xs text-muted-foreground">
            Downloads are preview-only and never create a file.
          </p>
        </div>
        <div className="relative w-full sm:max-w-xs">
          <Search className="pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search report..."
            aria-label="Search reports"
            className="rounded-full pl-9"
          />
        </div>
      </div>

      <div className="surface-card overflow-hidden hover:shadow-none">
        <div className="overflow-x-auto">
          <Table>
          <TableHeader>
            <TableRow className="hover:bg-transparent">
              <TableHead>Report ID</TableHead>
              <TableHead>Use case</TableHead>
              <TableHead className="min-w-40">Type</TableHead>
              <TableHead>Period</TableHead>
              <TableHead>Generated at</TableHead>
              <TableHead>Size</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="w-12">
                <span className="sr-only">Preview actions</span>
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {rows.map((report) => (
              <TableRow key={report.id}>
                <TableCell className="text-sm font-semibold">
                  {report.id}
                </TableCell>
                <TableCell className="text-sm">{report.useCase}</TableCell>
                <TableCell className="text-sm">{report.type}</TableCell>
                <TableCell className="text-xs text-muted-foreground">
                  {report.period}
                </TableCell>
                <TableCell className="text-xs whitespace-nowrap text-muted-foreground">
                  {report.generatedAt}
                </TableCell>
                <TableCell className="text-xs text-muted-foreground">
                  {report.size}
                </TableCell>
                <TableCell>
                  <ReportStatusBadge status={report.status} />
                </TableCell>
                <TableCell>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="rounded-full"
                    disabled={report.status !== "Ready"}
                    aria-label={`Preview download for ${report.id}`}
                    onClick={() =>
                      toast.info("Download preview simulated", {
                        description: `${report.id} was not downloaded. Preview only; no file was created.`,
                      })
                    }
                  >
                    <Download className="h-4 w-4" />
                  </Button>
                </TableCell>
              </TableRow>
            ))}
            {rows.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={8}
                  className="py-10 text-center text-sm text-muted-foreground"
                >
                  No report matches your search.
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
