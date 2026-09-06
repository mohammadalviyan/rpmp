import type { ExecutionTrend } from "@/lib/api/types";

const width = 760;
const height = 300;
const padding = { top: 24, right: 24, bottom: 44, left: 56 };
const plotWidth = width - padding.left - padding.right;
const plotHeight = height - padding.top - padding.bottom;
const tickCount = 4;

// Rounds the axis top to a readable step so gridline labels stay whole numbers.
function axisMaximum(value: number): number {
  if (value <= tickCount) {
    return tickCount;
  }

  const rough = value / tickCount;
  const magnitude = 10 ** Math.floor(Math.log10(rough));
  const normalized = rough / magnitude;
  const step =
    (normalized <= 1 ? 1 : normalized <= 2 ? 2 : normalized <= 5 ? 5 : 10) *
    magnitude;

  return step * tickCount;
}

function chartPoints(
  values: number[],
  maximum: number,
): Array<{ x: number; y: number }> {
  return values.map((value, index) => ({
    x:
      values.length === 1
        ? padding.left + plotWidth / 2
        : padding.left + (index / (values.length - 1)) * plotWidth,
    y: padding.top + plotHeight - (value / maximum) * plotHeight,
  }));
}

function toPolyline(points: Array<{ x: number; y: number }>): string {
  return points.map(({ x, y }) => `${x},${y}`).join(" ");
}

function toArea(points: Array<{ x: number; y: number }>): string {
  const baseline = padding.top + plotHeight;
  const line = points.map(({ x, y }) => `L ${x},${y}`).join(" ");
  const first = points[0];
  const last = points[points.length - 1];

  return `M ${first.x},${baseline} ${line} L ${last.x},${baseline} Z`;
}

export function ExecutionTrendChart({ trend }: { trend: ExecutionTrend }) {
  if (trend.points.length === 0) {
    return (
      <p role="status" className="mt-8 text-sm text-muted-foreground">
        No execution trend data is available for this period.
      </p>
    );
  }

  const maximum = axisMaximum(
    Math.max(
      1,
      ...trend.points.flatMap((point) => [point.success, point.failure]),
    ),
  );
  const successPoints = chartPoints(
    trend.points.map((point) => point.success),
    maximum,
  );
  const failurePoints = chartPoints(
    trend.points.map((point) => point.failure),
    maximum,
  );
  const ticks = Array.from(
    { length: tickCount + 1 },
    (_, index) => (maximum / tickCount) * index,
  );

  return (
    <>
      <div className="mt-5 flex gap-5 text-xs font-medium">
        <span className="flex items-center gap-2">
          <span className="h-2.5 w-2.5 rounded-full bg-primary" />
          Successful
        </span>
        <span className="flex items-center gap-2">
          <span className="h-2.5 w-2.5 rounded-full bg-destructive" />
          Failed
        </span>
      </div>

      <svg
        aria-label={`Execution trend with ${trend.points.length} time points. Successful executions are shown in blue and failed executions in red.`}
        className="mt-4 h-auto w-full overflow-visible"
        role="img"
        viewBox={`0 0 ${width} ${height}`}
      >
        <defs>
          <linearGradient id="trend-success-fill" x1="0" x2="0" y1="0" y2="1">
            <stop
              offset="0%"
              stopColor="var(--bri-primary-cakrawala-main)"
              stopOpacity="0.28"
            />
            <stop
              offset="100%"
              stopColor="var(--bri-primary-cakrawala-main)"
              stopOpacity="0.02"
            />
          </linearGradient>
          <linearGradient id="trend-failure-fill" x1="0" x2="0" y1="0" y2="1">
            <stop
              offset="0%"
              stopColor="var(--bri-red-main)"
              stopOpacity="0.24"
            />
            <stop offset="100%" stopColor="var(--bri-red-main)" stopOpacity="0.02" />
          </linearGradient>
        </defs>

        <g aria-hidden="true">
          {ticks.map((tick) => {
            const y =
              padding.top + plotHeight - (tick / maximum) * plotHeight;
            return (
              <g key={tick}>
                <line
                  stroke="var(--bri-black-300)"
                  strokeDasharray={tick === 0 ? undefined : "4 6"}
                  x1={padding.left}
                  x2={width - padding.right}
                  y1={y}
                  y2={y}
                />
                <text
                  dominantBaseline="middle"
                  fill="var(--bri-black-600)"
                  fontSize="12"
                  textAnchor="end"
                  x={padding.left - 12}
                  y={y}
                >
                  {tick.toLocaleString("en-US")}
                </text>
              </g>
            );
          })}

          <path d={toArea(successPoints)} fill="url(#trend-success-fill)" />
          <path d={toArea(failurePoints)} fill="url(#trend-failure-fill)" />

          <polyline
            fill="none"
            points={toPolyline(successPoints)}
            stroke="var(--bri-primary-cakrawala-main)"
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="3.5"
          />
          <polyline
            fill="none"
            points={toPolyline(failurePoints)}
            stroke="var(--bri-red-main)"
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="3.5"
          />

          {trend.points.map((point, index) => (
            <g key={`${point.bucket}-${point.label}`}>
              <circle
                cx={successPoints[index].x}
                cy={successPoints[index].y}
                fill="var(--bri-white-main)"
                r="5"
                stroke="var(--bri-primary-cakrawala-main)"
                strokeWidth="3"
              />
              <circle
                cx={failurePoints[index].x}
                cy={failurePoints[index].y}
                fill="var(--bri-white-main)"
                r="5"
                stroke="var(--bri-red-main)"
                strokeWidth="3"
              />
              <text
                fill="var(--bri-black-700)"
                fontSize="12"
                textAnchor="middle"
                x={successPoints[index].x}
                y={height - 14}
              >
                {point.label}
              </text>
            </g>
          ))}
        </g>
      </svg>

      <ul className="sr-only">
        {trend.points.map((point) => (
          <li key={point.bucket}>
            {point.label}: {point.success} successful, {point.failure} failed
          </li>
        ))}
      </ul>
    </>
  );
}
