import { redirect } from "next/navigation";

import { ApiOverviewDashboard } from "@/components/api-overview-dashboard";
import { OverviewDashboard } from "@/components/overview-dashboard";
import { getCurrentUser, getDashboardOverview } from "@/lib/api/server";
import {
  getRpaDataMode,
  getRpaDataProvider,
} from "@/lib/data/provider";

export default async function DashboardPage() {
  if (getRpaDataMode() === "api") {
    const [overview, user] = await Promise.all([
      getDashboardOverview(),
      getCurrentUser(),
    ]);

    if (overview.status === "unauthenticated") {
      redirect("/login");
    }

    return (
      <ApiOverviewDashboard
        canTriggerSync={user?.role === "admin"}
        errors={overview.errors}
        summary={overview.summary}
        trend={overview.trend}
      />
    );
  }

  const provider = getRpaDataProvider();
  const [
    summary,
    performanceSeries,
    errorDistribution,
    automationTypeSplit,
    useCases,
  ] = await Promise.all([
    provider.getSummary(),
    provider.getPerformanceSeries(),
    provider.getErrorDistribution(),
    provider.getAutomationTypeSplit(),
    provider.getUseCases(),
  ]);

  return (
    <OverviewDashboard
      automationTypeSplit={automationTypeSplit}
      errorDistribution={errorDistribution}
      performanceSeries={performanceSeries}
      summary={summary}
      useCases={useCases}
    />
  );
}
