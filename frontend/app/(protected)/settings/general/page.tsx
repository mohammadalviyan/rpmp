import type { Metadata } from "next";

import { GeneralSettingsPreview } from "@/components/settings-preview-pages";

export const metadata: Metadata = {
  title: "General Settings Preview | RPMP",
  description: "Preview workspace preferences for RPMP monitoring.",
};

export default function GeneralSettingsPage() {
  return <GeneralSettingsPreview />;
}
