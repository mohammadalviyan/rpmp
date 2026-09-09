import type { Metadata } from "next";

import { ReportHistoryPreview } from "@/components/report-history-preview";
import { getRpaDataProvider } from "@/lib/data/provider";

export const metadata: Metadata = {
  title: "Report History | RPMP",
  description:
    "Browse sample report history with preview-only download actions.",
};

export default async function ReportHistoryPage() {
  const reportHistory = await getRpaDataProvider().getReportHistory();

  return <ReportHistoryPreview reportHistory={reportHistory} />;
}
