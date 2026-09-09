export default function DashboardLoading() {
  return (
    <main
      aria-busy="true"
      aria-live="polite"
      className="mx-auto w-full max-w-[1600px] px-5 py-6 sm:px-8 sm:py-8"
    >
      <p className="text-sm text-muted-foreground">
        Loading current dashboard data...
      </p>
      <div
        aria-hidden="true"
        className="mt-5 grid gap-4 sm:grid-cols-2 xl:grid-cols-3"
      >
        {Array.from({ length: 6 }, (_, index) => (
          <div
            key={index}
            className="h-32 animate-pulse rounded-2xl border bg-card"
          />
        ))}
      </div>
      <div
        aria-hidden="true"
        className="mt-6 grid gap-6 lg:grid-cols-3"
      >
        <div className="h-96 animate-pulse rounded-2xl border bg-card lg:col-span-2" />
        <div className="h-96 animate-pulse rounded-2xl border bg-card" />
      </div>
    </main>
  );
}
