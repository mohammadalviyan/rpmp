import type { ReactNode } from "react";
import { redirect } from "next/navigation";

import {
  AppSidebar,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/app-sidebar";
import { LogoutControl } from "@/components/logout-control";
import { getCurrentUser } from "@/lib/api/server";

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
      <div className="min-h-screen bg-background lg:pl-64">
        <AppSidebar />

        <div className="min-w-0">
          <header className="sticky top-0 z-20 border-b bg-card/95 backdrop-blur">
            <div className="flex min-h-20 items-center justify-between gap-4 px-5 sm:px-8">
              <div className="flex min-w-0 items-center gap-3">
                <SidebarTrigger />
                <div className="min-w-0">
                  <h1 className="text-xl font-semibold">Overview</h1>
                  <p className="mt-0.5 truncate text-xs text-muted-foreground sm:text-sm">
                    Monitoring summary for your RPA operations
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-3 sm:gap-5">
                <div className="max-w-24 text-right sm:max-w-none">
                  <p className="truncate text-sm font-semibold">
                    {user.display_name}
                  </p>
                  <p className="text-xs capitalize text-muted-foreground">
                    {user.role}
                  </p>
                </div>
                <span
                  aria-hidden="true"
                  className="hidden h-10 w-10 place-items-center rounded-full bg-[var(--bri-primary-cakrawala-100)] text-sm font-bold text-[var(--bri-primary-cakrawala-700)] sm:grid"
                >
                  {user.display_name.slice(0, 1).toUpperCase()}
                </span>
                <LogoutControl />
              </div>
            </div>
          </header>
          {children}
        </div>
      </div>
    </SidebarProvider>
  );
}
