import { redirect } from "next/navigation";

import { DashboardSummary } from "@/components/dashboard-summary";
import { getDashboardSummary } from "@/lib/api/server";

export default async function DashboardPage() {
  const result = await getDashboardSummary();

  if (result.status === "unauthenticated") {
    redirect("/login");
  }

  return (
    <main className="mx-auto w-full max-w-6xl px-6 py-10">
      <h1 className="text-2xl font-semibold">Dashboard</h1>

      {result.status === "success" ? (
        <DashboardSummary summary={result.summary} />
      ) : (
        <section
          aria-labelledby="dashboard-unavailable-title"
          className="mt-6 rounded-xl border bg-card p-6"
        >
          <h2
            id="dashboard-unavailable-title"
            className="text-lg font-semibold"
          >
            Dashboard data unavailable
          </h2>
          <p className="mt-2 text-sm text-muted-foreground">
            {result.status === "source_unavailable"
              ? "Current dashboard data is temporarily unavailable. KPI values, period, and refresh time are not shown."
              : "We could not load current dashboard data. Try again later."}
          </p>
        </section>
      )}
    </main>
  );
}
