"use client";

import Link from "next/link";
import {
  AlertTriangle,
  ArrowUpRight,
  Bot,
  CheckCircle2,
  Layers,
  MailCheck,
  UserCheck,
} from "lucide-react";
import {
  Area,
  AreaChart,
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

import { StatusBadge, TypeBadge } from "@/components/rpa-badges";
import type {
  DistributionPoint,
  PerformancePoint,
  Summary,
  UseCase,
} from "@/lib/data/types";

type OverviewDashboardProps = {
  summary: Summary;
  performanceSeries: PerformancePoint[];
  errorDistribution: DistributionPoint[];
  automationTypeSplit: DistributionPoint[];
  useCases: UseCase[];
};

const cardClass =
  "rounded-2xl border bg-card shadow-[0_8px_24px_var(--bri-black-opacity-10)]";

const chartTooltip = {
  contentStyle: {
    background: "var(--color-card)",
    border: "1px solid var(--color-border)",
    borderRadius: "12px",
    color: "var(--color-foreground)",
    fontSize: "12px",
  },
};

export function OverviewDashboard({
  summary,
  performanceSeries,
  errorDistribution,
  automationTypeSplit,
  useCases,
}: OverviewDashboardProps) {
  const stats = [
    {
      label: "Total Use Cases",
      value: summary.totalUseCase,
      icon: Layers,
      hint: "+3 this month",
    },
    {
      label: "Success Rate",
      value: `${summary.successRate}%`,
      icon: CheckCircle2,
      hint: "+2.4% vs Aug",
    },
    {
      label: "Total Issues",
      value: summary.totalIssue,
      icon: AlertTriangle,
      hint: "5 open today",
    },
    {
      label: "Unattended Processes",
      value: summary.unattended,
      icon: Bot,
      hint: "24/7 scheduled",
    },
    {
      label: "Attended Processes",
      value: summary.attended,
      icon: UserCheck,
      hint: "agent assisted",
    },
    {
      label: "Reports Sent",
      value: summary.reportSent,
      icon: MailCheck,
      hint: "auto-delivered",
    },
  ];

  return (
    <main
      aria-label="Overview dashboard"
      className="mx-auto w-full max-w-[1600px] space-y-8 px-4 py-8 sm:px-6 lg:px-8"
    >
      <section
        aria-label="Dashboard key performance indicators"
        className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3"
      >
        {stats.map((stat) => (
          <article key={stat.label} className={`${cardClass} p-5`}>
            <div className="flex items-start justify-between gap-3">
              <div className="min-w-0">
                <h2 className="truncate text-xs font-medium tracking-wide text-muted-foreground uppercase">
                  {stat.label}
                </h2>
                <p className="mt-3 text-3xl font-bold tracking-tight tabular-nums">
                  {stat.value}
                </p>
                <p className="mt-1 text-xs text-muted-foreground">{stat.hint}</p>
              </div>
              <span className="grid h-10 w-10 shrink-0 place-items-center rounded-xl bg-primary/12 text-primary">
                <stat.icon aria-hidden="true" className="h-5 w-5" />
              </span>
            </div>
          </article>
        ))}
      </section>

      <section className="grid gap-4 lg:grid-cols-3">
        <article className={`${cardClass} min-w-0 p-5 lg:col-span-2`}>
          <div className="mb-4 flex flex-wrap items-end justify-between gap-2">
            <div>
              <h2 className="text-sm font-semibold">Automation Performance</h2>
              <p className="text-xs text-muted-foreground">
                Successful and failed processes per month
              </p>
            </div>
            <p className="text-xs text-muted-foreground">
              Success rate{" "}
              <span className="font-semibold text-primary">
                {summary.successRate}%
              </span>
            </p>
          </div>
          <div
            aria-label={`Automation performance chart with ${performanceSeries.length} monthly data points`}
            className="h-64"
            role="img"
          >
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={performanceSeries}>
                <defs>
                  <linearGradient id="performance-success" x1="0" x2="0" y1="0" y2="1">
                    <stop
                      offset="0%"
                      stopColor="var(--bri-primary-cakrawala-main)"
                      stopOpacity={0.5}
                    />
                    <stop
                      offset="100%"
                      stopColor="var(--bri-primary-cakrawala-main)"
                      stopOpacity={0}
                    />
                  </linearGradient>
                </defs>
                <CartesianGrid
                  stroke="var(--bri-black-300)"
                  strokeDasharray="3 3"
                  vertical={false}
                />
                <XAxis
                  dataKey="month"
                  stroke="var(--bri-black-300)"
                  tick={{
                    fill: "var(--bri-black-700)",
                    fontSize: 11,
                  }}
                />
                <YAxis
                  stroke="var(--bri-black-300)"
                  tick={{
                    fill: "var(--bri-black-700)",
                    fontSize: 11,
                  }}
                />
                <Tooltip {...chartTooltip} />
                <Area
                  dataKey="success"
                  fill="url(#performance-success)"
                  stroke="var(--bri-primary-cakrawala-main)"
                  strokeWidth={2}
                  type="monotone"
                />
                <Area
                  dataKey="failed"
                  fillOpacity={0}
                  stroke="var(--bri-red-main)"
                  strokeWidth={2}
                  type="monotone"
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
          <ul className="sr-only">
            {performanceSeries.map((point) => (
              <li key={point.month}>
                {point.month}: {point.success} successful, {point.failed} failed
              </li>
            ))}
          </ul>
        </article>

        <article className={`${cardClass} p-5`}>
          <h2 className="text-sm font-semibold">Error Distribution</h2>
          <p className="text-xs text-muted-foreground">
            Normalized error categories
          </p>
          <div
            aria-label={`Error distribution chart with ${errorDistribution.length} categories`}
            className="h-48"
            role="img"
          >
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={errorDistribution}
                  dataKey="value"
                  innerRadius={52}
                  nameKey="name"
                  outerRadius={76}
                  paddingAngle={3}
                  stroke="none"
                >
                  {errorDistribution.map((error, index) => (
                    <Cell
                      key={error.name}
                      fill={
                        index === 0
                          ? "var(--bri-red-main)"
                          : "var(--bri-yellow-main)"
                      }
                    />
                  ))}
                </Pie>
                <Tooltip {...chartTooltip} />
              </PieChart>
            </ResponsiveContainer>
          </div>
          <ul
            aria-label="Error distribution by category"
            className="space-y-2"
          >
            {errorDistribution.map((error, index) => (
              <li
                key={error.name}
                className="flex items-center justify-between text-xs"
              >
                <span className="flex items-center gap-2 text-muted-foreground">
                  <span
                    aria-hidden="true"
                    className="h-2 w-2 rounded-full"
                    style={{
                      backgroundColor:
                        index === 0
                          ? "var(--bri-red-main)"
                          : "var(--bri-yellow-main)",
                    }}
                  />
                  {error.name}
                </span>
                <span className="font-semibold">{error.value} issues</span>
              </li>
            ))}
          </ul>
        </article>
      </section>

      <section className="grid gap-4 lg:grid-cols-3">
        <article className={`${cardClass} flex h-[420px] flex-col p-5`}>
          <h2 className="text-sm font-semibold">Automation Type</h2>
          <p className="text-xs text-muted-foreground">
            Attended and unattended use cases
          </p>
          <div
            aria-label={`Automation type chart with ${automationTypeSplit.length} categories`}
            className="mt-2 min-h-0 w-full flex-1"
            role="img"
          >
            <ResponsiveContainer width="100%" height="100%">
              <BarChart
                barCategoryGap="25%"
                data={automationTypeSplit}
                layout="vertical"
                margin={{ top: 16, right: 24, bottom: 16, left: -16 }}
              >
                <CartesianGrid
                  horizontal={false}
                  stroke="var(--bri-black-300)"
                  strokeDasharray="4 4"
                />
                <XAxis
                  allowDecimals={false}
                  axisLine={{ stroke: "var(--bri-black-300)" }}
                  domain={[0, "dataMax + 2"]}
                  stroke="var(--bri-black-300)"
                  tick={{
                    fill: "var(--bri-black-700)",
                    fontSize: 13,
                    fontWeight: 500,
                  }}
                  tickLine={false}
                  type="number"
                />
                <YAxis
                  axisLine={false}
                  dataKey="name"
                  tick={{
                    fill: "var(--bri-black-main)",
                    fontSize: 14,
                    fontWeight: 600,
                  }}
                  tickLine={false}
                  type="category"
                  width={110}
                />
                <Tooltip
                  {...chartTooltip}
                  cursor={{ fill: "var(--bri-black-200)", opacity: 0.4 }}
                />
                <Bar dataKey="value" maxBarSize={64} radius={[0, 14, 14, 0]}>
                  {automationTypeSplit.map((type, index) => (
                    <Cell
                      key={type.name}
                      fill={
                        index === 0
                          ? "var(--bri-primary-cakrawala-main)"
                          : "var(--bri-secondary-mentari-main)"
                      }
                    />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>
          <ul className="sr-only">
            {automationTypeSplit.map((type) => (
              <li key={type.name}>
                {type.name}: {type.value} use cases
              </li>
            ))}
          </ul>
        </article>

        <article className={`${cardClass} p-5 lg:col-span-2`}>
          <div className="mb-4 flex flex-wrap items-center justify-between gap-2">
            <div>
              <h2 className="text-sm font-semibold">Top Use Cases</h2>
              <p className="text-xs text-muted-foreground">
                Success rate per automation
              </p>
            </div>
            <Link
              className="inline-flex items-center gap-1 text-xs font-semibold text-primary hover:underline"
              href="/use-cases"
            >
              View all <ArrowUpRight aria-hidden="true" className="h-3.5 w-3.5" />
            </Link>
          </div>
          <div className="space-y-4">
            {useCases.slice(0, 4).map((useCase) => (
              <Link
                key={useCase.id}
                className="block rounded-xl p-3 transition-colors hover:bg-muted/60"
                href={`/use-cases/${useCase.id}`}
              >
                <div className="grid grid-cols-[minmax(0,1fr)_auto] items-center gap-3">
                  <div className="min-w-0">
                    <p className="truncate text-sm font-semibold">
                      {useCase.name}
                    </p>
                    <p className="truncate text-xs text-muted-foreground">
                      {useCase.owner}
                    </p>
                  </div>
                  <div className="flex shrink-0 items-center gap-2">
                    <TypeBadge type={useCase.automationType} />
                    <StatusBadge status={useCase.status} />
                  </div>
                </div>
                <div className="mt-3 flex items-center gap-3">
                  <div
                    aria-label={`${useCase.successRate}% success rate`}
                    aria-valuemax={100}
                    aria-valuemin={0}
                    aria-valuenow={useCase.successRate}
                    className="h-1.5 flex-1 overflow-hidden rounded-full bg-muted"
                    role="progressbar"
                  >
                    <div
                      className="h-full bg-primary transition-all"
                      style={{ width: `${useCase.successRate}%` }}
                    />
                  </div>
                  <span className="w-10 shrink-0 text-right text-xs font-semibold">
                    {useCase.successRate}%
                  </span>
                </div>
              </Link>
            ))}
          </div>
        </article>
      </section>
    </main>
  );
}
