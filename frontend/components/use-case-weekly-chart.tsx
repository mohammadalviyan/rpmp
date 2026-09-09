"use client";

import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

import type { TrendPoint } from "@/lib/data/types";

export function UseCaseWeeklyChart({ trend }: { trend: TrendPoint[] }) {
  return (
    <>
      <div
        aria-label={`Weekly performance chart with ${trend.length} daily data points`}
        className="h-64"
        role="img"
      >
        <ResponsiveContainer height="100%" width="100%">
          <AreaChart
            data={trend}
            margin={{ top: 4, right: 8, left: -18, bottom: 0 }}
          >
            <defs>
              <linearGradient id="ucSuccess" x1="0" x2="0" y1="0" y2="1">
                <stop
                  offset="0%"
                  stopColor="var(--bri-green-main)"
                  stopOpacity={0.35}
                />
                <stop
                  offset="100%"
                  stopColor="var(--bri-green-main)"
                  stopOpacity={0}
                />
              </linearGradient>
              <linearGradient id="ucFailed" x1="0" x2="0" y1="0" y2="1">
                <stop
                  offset="0%"
                  stopColor="var(--bri-red-main)"
                  stopOpacity={0.35}
                />
                <stop
                  offset="100%"
                  stopColor="var(--bri-red-main)"
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
              axisLine={false}
              dataKey="day"
              tick={{
                fontSize: 11,
                fill: "var(--bri-black-700)",
              }}
              tickLine={false}
            />
            <YAxis
              axisLine={false}
              tick={{
                fontSize: 11,
                fill: "var(--bri-black-700)",
              }}
              tickLine={false}
            />
            <Tooltip
              contentStyle={{
                background: "var(--color-card)",
                border: "1px solid var(--color-border)",
                borderRadius: 12,
                fontSize: 12,
              }}
            />
            <Area
              dataKey="success"
              fill="url(#ucSuccess)"
              name="Success"
              stroke="var(--bri-green-main)"
              strokeWidth={2}
              type="monotone"
            />
            <Area
              dataKey="failed"
              fill="url(#ucFailed)"
              name="Failed"
              stroke="var(--bri-red-main)"
              strokeWidth={2}
              type="monotone"
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
      <ul className="sr-only">
        {trend.map((point) => (
          <li key={point.day}>
            {point.day}: {point.success} successful, {point.failed} failed
          </li>
        ))}
      </ul>
    </>
  );
}
