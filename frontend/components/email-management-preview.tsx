"use client";

import { Clock, Mail, Plus, Send, Trash2 } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { UseCase } from "@/lib/data/types";

type Schedule = {
  id: string;
  useCase: string;
  reportType: string;
  frequency: string;
  time: string;
  recipients: string[];
  active: boolean;
};

const initialSchedules: Schedule[] = [
  {
    id: "SCH-01",
    useCase: "Use Case A",
    reportType: "Performance Summary",
    frequency: "Monthly",
    time: "06:00",
    recipients: ["finance.ops@company.com", "management@company.com"],
    active: true,
  },
  {
    id: "SCH-02",
    useCase: "Use Case B",
    reportType: "Daily Operations",
    frequency: "Daily",
    time: "05:45",
    recipients: ["cs.lead@company.com"],
    active: true,
  },
  {
    id: "SCH-03",
    useCase: "Use Case E",
    reportType: "Compliance Report",
    frequency: "Weekly",
    time: "05:30",
    recipients: ["compliance@company.com"],
    active: false,
  },
];

type EmailManagementPreviewProps = {
  recipients: string[];
  reportTypes: string[];
  useCases: UseCase[];
};

export function EmailManagementPreview({
  recipients,
  reportTypes,
  useCases,
}: EmailManagementPreviewProps) {
  const [schedules, setSchedules] = useState(initialSchedules);
  const [useCase, setUseCase] = useState(useCases[0]?.name ?? "");
  const [reportType, setReportType] = useState(reportTypes[0] ?? "");
  const [frequency, setFrequency] = useState("Daily");
  const [time, setTime] = useState("06:00");
  const [selected, setSelected] = useState<string[]>([]);

  function toggleRecipient(email: string) {
    setSelected((current) =>
      current.includes(email)
        ? current.filter((recipient) => recipient !== email)
        : [...current, email],
    );
  }

  function addSchedule() {
    if (!useCase || !reportType) {
      toast.error("Schedule preview unavailable", {
        description: "Use case and report type data are required.",
      });
      return;
    }

    if (selected.length === 0) {
      toast.error("Pick at least one preview recipient", {
        description: "No schedule was created or saved.",
      });
      return;
    }

    setSchedules((current) => [
      ...current,
      {
        id: `SCH-${String(current.length + 1).padStart(2, "0")}`,
        useCase,
        reportType,
        frequency,
        time,
        recipients: selected,
        active: true,
      },
    ]);
    setSelected([]);
    toast.success("Schedule creation simulated", {
      description: `${reportType} for ${useCase} was added to this preview only. Nothing was saved and no email was sent.`,
    });
  }

  return (
    <main className="flex-1 px-4 py-8 sm:px-6 lg:px-8">
      <div className="mb-5 flex items-center gap-2">
        <span className="rounded-full border border-primary/30 bg-primary/10 px-2.5 py-1 text-xs font-semibold text-primary">
          Preview only
        </span>
        <p className="text-xs text-muted-foreground">
          Schedule actions stay in this browser view. Nothing is saved or sent.
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-[minmax(0,5fr)_minmax(0,4fr)]">
        <section className="space-y-4">
          <h2 className="text-sm font-bold">
            Email schedules ({schedules.length})
          </h2>
          {schedules.length === 0 ? (
            <div className="surface-card px-5 py-10 text-center hover:shadow-none">
              <p className="text-sm font-semibold">No preview schedules</p>
              <p className="mt-1 text-xs text-muted-foreground">
                Configure a schedule to add it to this browser view.
              </p>
            </div>
          ) : null}
          {schedules.map((schedule) => (
            <div
              key={schedule.id}
              className="surface-card p-5 hover:shadow-none"
            >
              <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div className="flex min-w-0 items-start gap-3">
                  <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-primary/12 text-primary">
                    <Mail className="h-4 w-4" />
                  </div>
                  <div className="min-w-0">
                    <p className="truncate text-sm font-bold">
                      {schedule.reportType} - {schedule.useCase}
                    </p>
                    <p className="mt-0.5 flex items-center gap-1.5 text-xs text-muted-foreground">
                      <Clock className="h-3 w-3" />
                      {schedule.frequency} at {schedule.time} ·{" "}
                      {schedule.recipients.length} recipient(s)
                    </p>
                    <div className="mt-2 flex flex-wrap gap-1.5">
                      {schedule.recipients.map((recipient) => (
                        <span
                          key={recipient}
                          className="rounded-full border border-border bg-muted px-2 py-0.5 text-[11px] text-muted-foreground"
                        >
                          {recipient}
                        </span>
                      ))}
                    </div>
                  </div>
                </div>
                <div className="flex shrink-0 items-center justify-between gap-1 sm:justify-start">
                  <span
                    className={
                      schedule.active
                        ? "rounded-full bg-success/12 px-2 py-0.5 text-[11px] font-semibold text-success"
                        : "rounded-full bg-muted px-2 py-0.5 text-[11px] font-semibold text-muted-foreground"
                    }
                  >
                    {schedule.active ? "Active" : "Paused"}
                  </span>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="rounded-full text-muted-foreground hover:text-destructive"
                    aria-label={`Remove preview schedule ${schedule.id}`}
                    onClick={() => {
                      setSchedules((current) =>
                        current.filter((item) => item.id !== schedule.id),
                      );
                      toast.success("Schedule removal simulated", {
                        description: `${schedule.id} was removed from this preview only. No saved schedule was changed.`,
                      });
                    }}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </div>
          ))}
        </section>

        <aside className="surface-card h-fit space-y-5 p-6 hover:shadow-none">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/12 text-primary">
              <Plus className="h-5 w-5" />
            </div>
            <div>
              <h2 className="text-sm font-bold">New email schedule</h2>
              <p className="text-xs text-muted-foreground">
                Configure a local schedule preview. Reports are not sent.
              </p>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="schedule-use-case">Use case</Label>
            <Select
              value={useCase}
              onValueChange={(value) => value && setUseCase(value)}
            >
              <SelectTrigger id="schedule-use-case" className="w-full">
                <SelectValue placeholder="Select use case" />
              </SelectTrigger>
              <SelectContent>
                {useCases.map((item) => (
                  <SelectItem key={item.id} value={item.name}>
                    {item.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="schedule-report-type">Report type</Label>
              <Select
                value={reportType}
                onValueChange={(value) => value && setReportType(value)}
              >
                <SelectTrigger id="schedule-report-type" className="w-full">
                  <SelectValue placeholder="Select report type" />
                </SelectTrigger>
                <SelectContent>
                  {reportTypes.map((type) => (
                    <SelectItem key={type} value={type}>
                      {type}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="schedule-frequency">Frequency</Label>
              <Select
                value={frequency}
                onValueChange={(value) => value && setFrequency(value)}
              >
                <SelectTrigger id="schedule-frequency" className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {["Daily", "Weekly", "Monthly"].map((value) => (
                    <SelectItem key={value} value={value}>
                      {value}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="schedule-time">Send time</Label>
            <Input
              id="schedule-time"
              type="time"
              value={time}
              onChange={(event) => setTime(event.target.value)}
            />
          </div>

          <div className="space-y-2">
            <Label>Recipients</Label>
            <div className="grid gap-2">
              {recipients.map((email) => (
                <label
                  key={email}
                  className="flex cursor-pointer items-center gap-2.5 rounded-lg border border-border bg-muted/40 px-3 py-2 text-xs transition-colors hover:border-primary/40"
                >
                  <Checkbox
                    checked={selected.includes(email)}
                    onCheckedChange={() => toggleRecipient(email)}
                  />
                  <span className="truncate">{email}</span>
                </label>
              ))}
            </div>
          </div>

          <Button
            className="w-full rounded-full"
            onClick={addSchedule}
            disabled={!useCase || !reportType}
          >
            <Send className="mr-2 h-4 w-4" />
            Simulate schedule
          </Button>
        </aside>
      </div>
    </main>
  );
}
