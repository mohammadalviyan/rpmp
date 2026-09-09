import type { Metadata } from "next";

import { GenerateReportPreview } from "@/components/generate-report-preview";
import { getRpaDataProvider } from "@/lib/data/provider";

export const metadata: Metadata = {
  title: "Generate Report | RPMP",
  description:
    "Preview report configuration and simulated delivery for monitored use cases.",
};

export default async function GenerateReportPage() {
  const provider = getRpaDataProvider();
  const [useCases, reportTypes, reportPeriods, recipients] = await Promise.all([
    provider.getUseCases(),
    provider.getReportTypes(),
    provider.getReportPeriods(),
    provider.getRecipients(),
  ]);

  return (
    <GenerateReportPreview
      recipients={recipients}
      reportPeriods={reportPeriods}
      reportTypes={reportTypes}
      useCases={useCases}
    />
  );
}
