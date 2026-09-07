import { Bell, Search } from "lucide-react";

import { LogoutControl } from "@/components/logout-control";
import { ThemeToggle } from "@/components/theme-toggle";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { SidebarTrigger } from "@/components/ui/sidebar";
import type { User } from "@/lib/api/types";

export function DashboardHeader({
  title,
  subtitle,
  user,
}: {
  title: string;
  subtitle?: string;
  user: User;
}) {
  return (
    <header className="sticky top-0 z-30 border-b border-border bg-background/85 backdrop-blur-xl">
      <div className="grid grid-cols-[minmax(0,1fr)_auto] items-center gap-3 px-4 py-3 sm:px-6 lg:px-8">
        <div className="flex min-w-0 items-center gap-3">
          <SidebarTrigger className="shrink-0" />
          <div className="min-w-0">
            <h1 className="truncate text-base font-bold tracking-tight sm:text-lg">
              {title}
            </h1>
            {subtitle ? (
              <p className="hidden truncate text-xs text-muted-foreground sm:block">
                {subtitle}
              </p>
            ) : null}
          </div>
        </div>

        <div className="flex items-center gap-2">
          <div className="relative hidden md:block">
            <Search className="pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              aria-label="Search"
              placeholder="Search use case, report, issue..."
              className="w-56 rounded-full pl-9 lg:w-72"
            />
          </div>
          <Button
            variant="ghost"
            size="icon"
            className="relative rounded-full"
            aria-label="Notifications"
          >
            <Bell className="h-4 w-4" />
            <span className="absolute top-2 right-2 h-2 w-2 rounded-full bg-primary" />
          </Button>
          <ThemeToggle />
          <div className="flex min-w-0 items-center gap-2 pl-1">
            <Avatar className="h-8 w-8 shrink-0">
              <AvatarFallback className="bg-[var(--bri-primary-cakrawala-main)] text-xs font-bold text-white">
                {user.display_name.slice(0, 1).toUpperCase()}
              </AvatarFallback>
            </Avatar>
            <div className="hidden min-w-0 lg:block">
              <p className="truncate text-xs font-semibold">
                {user.display_name}
              </p>
              <p className="truncate text-[11px] capitalize text-muted-foreground">
                {user.role}
              </p>
            </div>
          </div>
          <LogoutControl />
        </div>
      </div>
    </header>
  );
}
