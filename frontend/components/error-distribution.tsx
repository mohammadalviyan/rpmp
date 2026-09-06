import type { DashboardErrors } from "@/lib/api/types";

// BRI tokens only. Failure categories stay in the red and amber range so the
// donut reads as failure data, not as a success accent.
const segmentColors = [
  "var(--bri-red-main)",
  "var(--bri-yellow-main)",
  "var(--bri-black-500)",
  "var(--bri-red-700)",
  "var(--bri-yellow-600)",
];

const size = 168;
const center = size / 2;
const radius = 64;
const ringWidth = 26;
const segmentGap = 1.2;

export function ErrorDistribution({
  errors,
}: {
  errors: DashboardErrors;
}) {
  const groups = errors.groups.filter((group) => group.count > 0);

  if (groups.length === 0) {
    return (
      <p role="status" className="mt-8 text-sm text-muted-foreground">
        No errors were recorded for this period.
      </p>
    );
  }

  const total = groups.reduce((sum, group) => sum + group.count, 0);
  const segments = groups.map((group, index) => {
    const preceding = groups
      .slice(0, index)
      .reduce((sum, earlier) => sum + earlier.count, 0);

    return {
      ...group,
      color: segmentColors[index % segmentColors.length],
      percentage: (group.count / total) * 100,
      offset: (preceding / total) * 100,
    };
  });

  return (
    <>
      <div className="mt-6 flex justify-center">
        <svg
          aria-label={`Error distribution donut chart. ${total} failed executions across ${groups.length} categories.`}
          className="h-auto w-40"
          role="img"
          viewBox={`0 0 ${size} ${size}`}
        >
          <g transform={`rotate(-90 ${center} ${center})`}>
            <circle
              cx={center}
              cy={center}
              fill="none"
              r={radius}
              stroke="var(--bri-black-200)"
              strokeWidth={ringWidth}
            />
            {segments.map((segment) => {
              const length = Math.max(
                segment.percentage - (segments.length > 1 ? segmentGap : 0),
                0.4,
              );

              return (
                <circle
                  key={segment.code}
                  cx={center}
                  cy={center}
                  fill="none"
                  pathLength={100}
                  r={radius}
                  stroke={segment.color}
                  strokeDasharray={`${length} ${100 - length}`}
                  strokeDashoffset={-segment.offset}
                  strokeWidth={ringWidth}
                />
              );
            })}
          </g>
          <text
            fill="var(--bri-black-main)"
            fontSize="30"
            fontWeight="600"
            textAnchor="middle"
            x={center}
            y={center + 2}
          >
            {total.toLocaleString("en-US")}
          </text>
          <text
            fill="var(--bri-black-600)"
            fontSize="12"
            textAnchor="middle"
            x={center}
            y={center + 22}
          >
            failures
          </text>
        </svg>
      </div>

      <ul
        aria-label="Error distribution by category"
        className="mt-6 space-y-3"
      >
        {segments.map((segment) => (
          <li
            key={segment.code}
            className="flex items-center justify-between gap-3 text-sm"
          >
            <span className="flex min-w-0 items-center gap-2.5">
              <span
                aria-hidden="true"
                className="size-2.5 shrink-0 rounded-full"
                style={{ backgroundColor: segment.color }}
              />
              <span className="truncate font-medium">{segment.label}</span>
            </span>
            <span className="shrink-0 tabular-nums">
              <span className="font-semibold">
                {segment.count.toLocaleString("en-US")}
              </span>
              <span className="ml-2 text-muted-foreground">
                {Math.round(segment.percentage)}%
              </span>
            </span>
          </li>
        ))}
      </ul>
    </>
  );
}
