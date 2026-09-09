"use client";

import { usePathname } from "next/navigation";

import { DashboardHeader } from "@/components/dashboard-header";
import type { User } from "@/lib/api/types";

const routeHeaders = {
  dashboard: {
    title: "Overview",
    subtitle: "Monitoring summary for your RPA operations",
  },
  useCases: {
    title: "Use Cases",
    subtitle: "All monitored automation processes",
  },
  useCaseDetail: {
    title: "Use Case Details",
    subtitle: "Performance and issue history for one automation process",
  },
  generateReport: {
    title: "Generate Report",
    subtitle: "Configure a preview-only report simulation",
  },
  reportHistory: {
    title: "Report History",
    subtitle: "Sample report records with preview-only downloads",
  },
  emailManagement: {
    title: "Email Management",
    subtitle: "Preview report delivery schedules without sending email",
  },
  emailHistory: {
    title: "Email History",
    subtitle: "Sample report email delivery records",
  },
  generalSettings: {
    title: "General Settings",
    subtitle: "Preview workspace preferences without saving changes",
  },
  userManagement: {
    title: "User Management",
    subtitle: "Preview sample members and roles without changing access",
  },
  emailSettings: {
    title: "Email Configuration",
    subtitle: "Preview email settings without connecting or sending",
  },
} as const;

export function ProtectedHeader({ user }: { user: User }) {
  const pathname = usePathname();
  const header =
    pathname === "/use-cases"
      ? routeHeaders.useCases
      : pathname.startsWith("/use-cases/")
        ? routeHeaders.useCaseDetail
        : pathname === "/generate-report"
          ? routeHeaders.generateReport
          : pathname === "/report-history"
            ? routeHeaders.reportHistory
            : pathname === "/email-management"
              ? routeHeaders.emailManagement
              : pathname === "/email-history"
                ? routeHeaders.emailHistory
                : pathname === "/settings/general"
                  ? routeHeaders.generalSettings
                  : pathname === "/settings/users"
                    ? routeHeaders.userManagement
                    : pathname === "/settings/email"
                      ? routeHeaders.emailSettings
                      : routeHeaders.dashboard;

  return <DashboardHeader {...header} user={user} />;
}
