"use client";

import { RotateCw, Search } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

import { EmailStatusBadge } from "@/components/rpa-badges";
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
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import type { EmailRecord, EmailStatus } from "@/lib/data/types";

type StatusFilter = "all" | Lowercase<EmailStatus>;

export function EmailHistoryPreview({
  emailHistory,
}: {
  emailHistory: EmailRecord[];
}) {
  const [query, setQuery] = useState("");
  const [filter, setFilter] = useState<StatusFilter>("all");

  const rows = emailHistory.filter((email) => {
    const matchesQuery = `${email.recipient}${email.subject}${email.useCase}`
      .toLowerCase()
      .includes(query.toLowerCase());
    const matchesStatus =
      filter === "all" || email.status.toLowerCase() === filter;
    return matchesQuery && matchesStatus;
  });

  return (
    <main className="flex-1 space-y-6 px-4 py-8 sm:px-6 lg:px-8">
      <div className="flex items-center gap-2">
        <span className="rounded-full border border-primary/30 bg-primary/10 px-2.5 py-1 text-xs font-semibold text-primary">
          Demo data
        </span>
        <p className="text-xs text-muted-foreground">
          Resend actions are local previews. No email is queued or sent.
        </p>
      </div>

      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="relative min-w-0">
          <Search className="pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search email..."
            aria-label="Search email history"
            className="w-full rounded-full pl-9 sm:w-72"
          />
        </div>
        <Tabs
          className="max-w-full overflow-x-auto"
          value={filter}
          onValueChange={(value) => setFilter(value as StatusFilter)}
        >
          <TabsList className="rounded-full">
            <TabsTrigger value="all" className="rounded-full text-xs">
              All
            </TabsTrigger>
            <TabsTrigger value="sent" className="rounded-full text-xs">
              Sent
            </TabsTrigger>
            <TabsTrigger value="pending" className="rounded-full text-xs">
              Pending
            </TabsTrigger>
            <TabsTrigger value="failed" className="rounded-full text-xs">
              Failed
            </TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      <div className="surface-card overflow-hidden hover:shadow-none">
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow className="hover:bg-transparent">
                <TableHead>Email ID</TableHead>
                <TableHead>Use case</TableHead>
                <TableHead className="min-w-48">Recipient</TableHead>
                <TableHead className="min-w-56">Subject</TableHead>
                <TableHead>Date</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="w-12">
                  <span className="sr-only">Preview actions</span>
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {rows.map((email) => (
                <TableRow key={email.id}>
                  <TableCell className="text-sm font-semibold">
                    {email.id}
                  </TableCell>
                  <TableCell className="text-sm">{email.useCase}</TableCell>
                  <TableCell className="text-sm">{email.recipient}</TableCell>
                  <TableCell className="text-xs text-muted-foreground">
                    {email.subject}
                  </TableCell>
                  <TableCell className="text-xs whitespace-nowrap text-muted-foreground">
                    {email.date}
                  </TableCell>
                  <TableCell>
                    <EmailStatusBadge status={email.status} />
                  </TableCell>
                  <TableCell>
                    {email.status === "Failed" ? (
                      <Button
                        variant="ghost"
                        size="icon"
                        className="rounded-full"
                        aria-label={`Preview resend for ${email.id}`}
                        onClick={() =>
                          toast.success("Resend preview simulated", {
                            description: `${email.id} was not queued or sent to ${email.recipient}.`,
                          })
                        }
                      >
                        <RotateCw className="h-4 w-4" />
                      </Button>
                    ) : null}
                  </TableCell>
                </TableRow>
              ))}
              {rows.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={7}
                    className="py-10 text-center text-sm text-muted-foreground"
                  >
                    No email matches your filters.
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
