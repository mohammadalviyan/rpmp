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
} as const;

export function ProtectedHeader({ user }: { user: User }) {
  const pathname = usePathname();
  const header =
    pathname === "/use-cases"
      ? routeHeaders.useCases
      : pathname.startsWith("/use-cases/")
        ? routeHeaders.useCaseDetail
        : routeHeaders.dashboard;

  return <DashboardHeader {...header} user={user} />;
}
