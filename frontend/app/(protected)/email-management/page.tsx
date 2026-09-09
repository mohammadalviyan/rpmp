import type { Metadata } from "next";

import { EmailManagementPreview } from "@/components/email-management-preview";
import { getRpaDataProvider } from "@/lib/data/provider";

export const metadata: Metadata = {
  title: "Email Management | RPMP",
  description:
    "Preview report email schedules and recipients for monitored use cases.",
};

export default async function EmailManagementPage() {
  const provider = getRpaDataProvider();
  const [useCases, reportTypes, recipients] = await Promise.all([
    provider.getUseCases(),
    provider.getReportTypes(),
    provider.getRecipients(),
  ]);

  return (
    <EmailManagementPreview
      recipients={recipients}
      reportTypes={reportTypes}
      useCases={useCases}
    />
  );
}
