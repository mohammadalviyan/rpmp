import type {
  AutomationType,
  EmailStatus,
  ErrorType,
  IssueStatus,
  UseCaseStatus,
} from "@/lib/data/types";
import { cn } from "@/lib/utils";

const base =
  "inline-flex items-center gap-1.5 whitespace-nowrap rounded-full border px-2.5 py-0.5 text-[11px] font-semibold";

export function TypeBadge({ type }: { type: AutomationType }) {
  return (
    <span
      className={cn(
        base,
        type === "Unattended"
          ? "border-primary/35 bg-primary/12 text-primary"
          : "border-border bg-muted text-muted-foreground",
      )}
    >
      <span
        className={cn(
          "h-1.5 w-1.5 rounded-full",
          type === "Unattended" ? "bg-primary" : "bg-muted-foreground",
        )}
      />
      {type}
    </span>
  );
}

export function StatusBadge({ status }: { status: UseCaseStatus }) {
  const classes: Record<UseCaseStatus, string> = {
    Running:
      "border-[var(--bri-green-300)] bg-[var(--bri-green-100)] text-[var(--bri-green-700)]",
    Warning:
      "border-[var(--bri-yellow-300)] bg-[var(--bri-yellow-100)] text-[var(--bri-yellow-700)]",
    Stopped:
      "border-[var(--bri-red-300)] bg-[var(--bri-red-100)] text-[var(--bri-red-700)]",
  };

  return <span className={cn(base, classes[status])}>{status}</span>;
}

export function ErrorTypeBadge({ type }: { type: ErrorType }) {
  return (
    <span
      className={cn(
        base,
        type === "Error A"
          ? "border-[var(--bri-red-300)] bg-[var(--bri-red-100)] text-[var(--bri-red-700)]"
          : "border-[var(--bri-yellow-300)] bg-[var(--bri-yellow-100)] text-[var(--bri-yellow-700)]",
      )}
    >
      {type}
    </span>
  );
}

export function IssueStatusBadge({ status }: { status: IssueStatus }) {
  const classes: Record<IssueStatus, string> = {
    Open:
      "border-[var(--bri-red-300)] bg-[var(--bri-red-100)] text-[var(--bri-red-700)]",
    "In Progress":
      "border-[var(--bri-yellow-300)] bg-[var(--bri-yellow-100)] text-[var(--bri-yellow-700)]",
    Resolved:
      "border-[var(--bri-green-300)] bg-[var(--bri-green-100)] text-[var(--bri-green-700)]",
  };

  return <span className={cn(base, classes[status])}>{status}</span>;
}

export function EmailStatusBadge({ status }: { status: EmailStatus }) {
  const classes: Record<EmailStatus, string> = {
    Sent:
      "border-[var(--bri-green-300)] bg-[var(--bri-green-100)] text-[var(--bri-green-700)]",
    Pending:
      "border-[var(--bri-yellow-300)] bg-[var(--bri-yellow-100)] text-[var(--bri-yellow-700)]",
    Failed:
      "border-[var(--bri-red-300)] bg-[var(--bri-red-100)] text-[var(--bri-red-700)]",
  };

  return <span className={cn(base, classes[status])}>{status}</span>;
}

export function ReportStatusBadge({
  status,
}: {
  status: "Ready" | "Generating" | "Failed";
}) {
  const classes = {
    Ready:
      "border-[var(--bri-green-300)] bg-[var(--bri-green-100)] text-[var(--bri-green-700)]",
    Generating:
      "border-[var(--bri-yellow-300)] bg-[var(--bri-yellow-100)] text-[var(--bri-yellow-700)]",
    Failed:
      "border-[var(--bri-red-300)] bg-[var(--bri-red-100)] text-[var(--bri-red-700)]",
  };

  return <span className={cn(base, classes[status])}>{status}</span>;
}
