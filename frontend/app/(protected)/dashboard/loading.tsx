export default function DashboardLoading() {
  return (
    <main
      aria-busy="true"
      aria-live="polite"
      className="mx-auto w-full max-w-6xl px-6 py-10"
    >
      <h1 className="text-2xl font-semibold">Dashboard</h1>
      <p className="mt-2 text-sm text-muted-foreground">
        Loading current dashboard data...
      </p>
      <div
        aria-hidden="true"
        className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-5"
      >
        {Array.from({ length: 5 }, (_, index) => (
          <div
            key={index}
            className="h-32 animate-pulse rounded-xl border bg-card"
          />
        ))}
      </div>
    </main>
  );
}
