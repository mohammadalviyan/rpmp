"use client";

import { FileBarChart, Loader2, Mail, Send, Sparkles } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { UseCase } from "@/lib/data/types";

type GenerateReportPreviewProps = {
  recipients: string[];
  reportPeriods: string[];
  reportTypes: string[];
  useCases: UseCase[];
};

export function GenerateReportPreview({
  recipients,
  reportPeriods,
  reportTypes,
  useCases,
}: GenerateReportPreviewProps) {
  const [useCaseId, setUseCaseId] = useState(useCases[0]?.id ?? "");
  const [reportType, setReportType] = useState(reportTypes[0] ?? "");
  const [period, setPeriod] = useState(reportPeriods[3] ?? reportPeriods[0] ?? "");
  const [selected, setSelected] = useState<string[]>(
    recipients[0] ? [recipients[0]] : [],
  );
  const [previewEmail, setPreviewEmail] = useState(true);
  const [generating, setGenerating] = useState(false);

  const active = useCases.find((useCase) => useCase.id === useCaseId);

  function toggleRecipient(email: string) {
    setSelected((current) =>
      current.includes(email)
        ? current.filter((recipient) => recipient !== email)
        : [...current, email],
    );
  }

  function simulateGenerate() {
    if (!active) {
      toast.error("Report preview unavailable", {
        description: "No use case data is available for this simulation.",
      });
      return;
    }

    if (selected.length === 0 && previewEmail) {
      toast.error("Select a preview recipient", {
        description:
          "Choose at least one sample recipient or disable the email-delivery preview.",
      });
      return;
    }

    setGenerating(true);
    window.setTimeout(() => {
      setGenerating(false);
      toast.success("Report generation simulated", {
        description: previewEmail
          ? `${reportType} for ${active.name} (${period}) was previewed for ${selected.length} recipient(s). No email was sent and no history was changed.`
          : `${reportType} for ${active.name} (${period}) was previewed only. No file was created and no history was changed.`,
      });
    }, 800);
  }

  if (!active || reportTypes.length === 0 || reportPeriods.length === 0) {
    return (
      <main className="flex-1 px-4 py-8 sm:px-6 lg:px-8">
        <div className="surface-card p-6 text-sm text-muted-foreground">
          Report preview data is unavailable.
        </div>
      </main>
    );
  }

  return (
    <main className="flex-1 px-4 py-8 sm:px-6 lg:px-8">
      <div className="mb-5 flex items-center gap-2">
        <span className="rounded-full border border-primary/30 bg-primary/10 px-2.5 py-1 text-xs font-semibold text-primary">
          Preview only
        </span>
        <p className="text-xs text-muted-foreground">
          Actions on this page are simulations. Nothing is saved, sent, or downloaded.
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-[minmax(0,5fr)_minmax(0,4fr)]">
        <section className="surface-card space-y-6 p-6 hover:shadow-none">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/12 text-primary">
              <FileBarChart className="h-5 w-5" />
            </div>
            <div>
              <h2 className="text-sm font-bold">Report configuration</h2>
              <p className="text-xs text-muted-foreground">
                Pick a use case, format, and period for the preview.
              </p>
            </div>
          </div>

          <div className="grid gap-5 sm:grid-cols-2">
            <div className="space-y-2 sm:col-span-2">
              <Label htmlFor="report-use-case">Use case</Label>
              <Select
                value={useCaseId}
                onValueChange={(value) => value && setUseCaseId(value)}
              >
                <SelectTrigger id="report-use-case" className="w-full">
                  <SelectValue placeholder="Select use case" />
                </SelectTrigger>
                <SelectContent>
                  {useCases.map((useCase) => (
                    <SelectItem key={useCase.id} value={useCase.id}>
                      {useCase.name} - {useCase.owner}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="report-type">Report type</Label>
              <Select
                value={reportType}
                onValueChange={(value) => value && setReportType(value)}
              >
                <SelectTrigger id="report-type" className="w-full">
                  <SelectValue placeholder="Select type" />
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
              <Label htmlFor="report-period">Period</Label>
              <Select
                value={period}
                onValueChange={(value) => value && setPeriod(value)}
              >
                <SelectTrigger id="report-period" className="w-full">
                  <SelectValue placeholder="Select period" />
                </SelectTrigger>
                <SelectContent>
                  {reportPeriods.map((reportPeriod) => (
                    <SelectItem key={reportPeriod} value={reportPeriod}>
                      {reportPeriod}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-3">
            <div className="flex items-center justify-between gap-4">
              <Label>Email recipients</Label>
              <label className="flex cursor-pointer items-center gap-2 text-xs text-muted-foreground">
                <Checkbox
                  checked={previewEmail}
                  onCheckedChange={setPreviewEmail}
                />
                Preview email delivery only
              </label>
            </div>
            <div
              className={
                previewEmail
                  ? "grid gap-2 sm:grid-cols-2"
                  : "grid gap-2 opacity-45 sm:grid-cols-2"
              }
            >
              {recipients.map((email) => (
                <label
                  key={email}
                  className="flex cursor-pointer items-center gap-2.5 rounded-lg border border-border bg-muted/40 px-3 py-2.5 text-xs transition-colors hover:border-primary/40"
                >
                  <Checkbox
                    disabled={!previewEmail}
                    checked={selected.includes(email)}
                    onCheckedChange={() => toggleRecipient(email)}
                  />
                  <span className="truncate">{email}</span>
                </label>
              ))}
            </div>
          </div>

          <div className="flex flex-col gap-2 border-t border-border pt-5 sm:flex-row sm:justify-end">
            <Button
              variant="outline"
              className="rounded-full"
              disabled={generating}
              onClick={() =>
                toast.info("Draft save simulated", {
                  description:
                    "Preview only. No draft was saved and report history was not changed.",
                })
              }
            >
              <Sparkles className="mr-2 h-4 w-4" />
              Simulate draft save
            </Button>
            <Button
              className="rounded-full"
              onClick={simulateGenerate}
              disabled={generating}
            >
              {generating ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <Send className="mr-2 h-4 w-4" />
              )}
              {generating ? "Simulating..." : "Simulate report"}
            </Button>
          </div>
        </section>

        <aside className="space-y-6">
          <div className="surface-card p-6 hover:shadow-none">
            <p className="text-[11px] font-semibold tracking-[0.18em] text-muted-foreground uppercase">
              Report preview
            </p>
            <h3 className="mt-3 text-lg font-bold">{reportType}</h3>
            <p className="text-sm text-muted-foreground">
              {active.name} · {period}
            </p>

            <dl className="mt-5 space-y-3 border-t border-border pt-5 text-sm">
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Owner</dt>
                <dd className="font-semibold">{active.owner}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Success rate</dt>
                <dd className="font-semibold text-success">{active.successRate}%</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Total processed</dt>
                <dd className="font-semibold">
                  {active.totalProcess.toLocaleString()}
                </dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Failed</dt>
                <dd className="font-semibold text-destructive">
                  {active.totalFailed}
                </dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">Open issues</dt>
                <dd className="font-semibold">
                  {
                    active.issues.filter((issue) => issue.status === "Open")
                      .length
                  }
                </dd>
              </div>
            </dl>
          </div>

          <div className="surface-card flex items-start gap-3 p-5 hover:shadow-none">
            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-primary/12 text-primary">
              <Mail className="h-4 w-4" />
            </div>
            <p className="text-xs leading-relaxed text-muted-foreground">
              Email delivery is preview-only. The simulation does not send mail,
              create a report file, or add an entry to report history.
            </p>
          </div>
        </aside>
      </div>
    </main>
  );
}
