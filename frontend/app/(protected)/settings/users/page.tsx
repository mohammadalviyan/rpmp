import type { Metadata } from "next";

import { UserManagementPreview } from "@/components/settings-preview-pages";

export const metadata: Metadata = {
  title: "User Management Preview | RPMP",
  description: "Preview sample RPMP members and roles without changing access.",
};

export default function UserManagementPage() {
  return <UserManagementPreview />;
}
