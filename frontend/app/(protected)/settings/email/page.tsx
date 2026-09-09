import type { Metadata } from "next";

import { EmailSettingsPreview } from "@/components/settings-preview-pages";

export const metadata: Metadata = {
  title: "Email Configuration Preview | RPMP",
  description: "Preview local email settings without connecting or sending.",
};

export default function EmailSettingsPage() {
  return <EmailSettingsPreview />;
}
