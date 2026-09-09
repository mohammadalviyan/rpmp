import type { ReactNode } from "react";
import { redirect } from "next/navigation";

import { AppSidebar } from "@/components/app-sidebar";
import { DemoDataBanner } from "@/components/demo-data-banner";
import { ProtectedHeader } from "@/components/protected-header";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { getCurrentUser } from "@/lib/api/server";
import { getRpaDataMode } from "@/lib/data/provider";

export default async function ProtectedLayout({
  children,
}: {
  children: ReactNode;
}) {
  const user = await getCurrentUser();

  if (!user) {
    redirect("/login");
  }

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <ProtectedHeader user={user} />
        {getRpaDataMode() === "mock" ? <DemoDataBanner /> : null}
        {children}
      </SidebarInset>
    </SidebarProvider>
  );
}
