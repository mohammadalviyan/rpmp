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
        className="mt-5 grid gap-4 sm:grid-cols-2 xl:grid-cols-5"
      >
        {Array.from({ length: 5 }, (_, index) => (
          <div
            key={index}
            className="h-32 animate-pulse rounded-2xl border bg-card"
          />
        ))}
      </div>
      <div
        aria-hidden="true"
        className="mt-6 grid gap-6 xl:grid-cols-[minmax(0,2fr)_minmax(19rem,1fr)]"
      >
        <div className="h-96 animate-pulse rounded-2xl border bg-card" />
        <div className="h-96 animate-pulse rounded-2xl border bg-card" />
      </div>
    </main>
  );
}
