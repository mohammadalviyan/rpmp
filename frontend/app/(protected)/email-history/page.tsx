import type { Metadata } from "next";

import { EmailHistoryPreview } from "@/components/email-history-preview";
import { getRpaDataProvider } from "@/lib/data/provider";

export const metadata: Metadata = {
  title: "Email History | RPMP",
  description: "Preview the report email delivery log for monitored use cases.",
};

export default async function EmailHistoryPage() {
  const emailHistory = await getRpaDataProvider().getEmailHistory();

  return <EmailHistoryPreview emailHistory={emailHistory} />;
}
