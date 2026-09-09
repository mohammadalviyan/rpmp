import { OverviewDashboard } from "@/components/overview-dashboard";
import { getRpaDataProvider } from "@/lib/data/provider";

export default async function DashboardPage() {
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
