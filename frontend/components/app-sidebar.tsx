"use client";

import {
  Boxes,
  FileChartColumn,
  History,
  LayoutDashboard,
  Mail,
  MailCheck,
  ServerCog,
  SlidersHorizontal,
  Users,
} from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar";

const groups = [
  {
    label: "Dashboard",
    items: [
      { title: "Overview", href: "/dashboard", icon: LayoutDashboard },
    ],
  },
  {
    label: "Report Management",
    items: [
      { title: "Use Cases", href: "/use-cases", icon: Boxes },
      {
        title: "Generate Report",
        href: "/generate-report",
        icon: FileChartColumn,
      },
      { title: "Report History", href: "/report-history", icon: History },
    ],
  },
  {
    label: "Email",
    items: [
      { title: "Email Management", href: "/email-management", icon: Mail },
      { title: "Email History", href: "/email-history", icon: MailCheck },
    ],
  },
  {
    label: "Settings",
    items: [
      {
        title: "General Settings",
        href: "/settings/general",
        icon: SlidersHorizontal,
      },
      { title: "User Management", href: "/settings/users", icon: Users },
      {
        title: "Email Configuration",
        href: "/settings/email",
        icon: ServerCog,
      },
    ],
  },
] as const;

export function AppSidebar() {
  const { state, setOpenMobile } = useSidebar();
  const collapsed = state === "collapsed";
  const pathname = usePathname();
  const isActive = (href: string) =>
    href === "/dashboard"
      ? pathname === href
      : pathname === href || pathname.startsWith(`${href}/`);

  return (
    <Sidebar collapsible="icon" className="border-sidebar-border">
      <SidebarHeader className="px-3 py-4">
        <div className="flex min-w-0 items-center gap-3">
          <Link
            aria-label="RPMP Overview"
            className="flex min-w-0 items-center gap-3"
            href="/dashboard"
            onClick={() => setOpenMobile(false)}
          >
            <span className="grid h-9 w-9 shrink-0 place-items-center rounded-lg bg-primary text-sm font-bold text-primary-foreground">
              R
            </span>
            {!collapsed ? (
              <span className="min-w-0">
                <span className="block truncate text-sm font-bold tracking-[0.18em]">
                  RPMP
                </span>
                <span className="block truncate text-[11px] text-sidebar-foreground/60">
                  Performance Monitoring
                </span>
              </span>
            ) : null}
          </Link>
        </div>
      </SidebarHeader>

      <SidebarContent className="gap-1">
        <nav aria-label="Primary navigation">
          {groups.map((group) => (
            <SidebarGroup key={group.label}>
              <SidebarGroupLabel className="text-[10px] font-semibold uppercase tracking-[0.2em] text-sidebar-foreground/45">
                {group.label}
              </SidebarGroupLabel>
              <SidebarGroupContent>
                <SidebarMenu>
                  {group.items.map((item) => {
                    const active = isActive(item.href);
                    return (
                      <SidebarMenuItem key={item.title}>
                        <SidebarMenuButton
                          asChild
                          isActive={active}
                          tooltip={item.title}
                          className="data-[active=true]:bg-[var(--bri-primary-cakrawala-100)] data-[active=true]:font-semibold data-[active=true]:text-[var(--bri-primary-cakrawala-700)]"
                        >
                          <Link
                            aria-current={active ? "page" : undefined}
                            className="flex items-center gap-3"
                            href={item.href}
                            onClick={() => setOpenMobile(false)}
                          >
                            <item.icon
                              aria-hidden="true"
                              className={
                                active
                                  ? "size-4 shrink-0 text-[var(--bri-primary-cakrawala-main)]"
                                  : "size-4 shrink-0"
                              }
                            />
                            <span className="truncate">{item.title}</span>
                          </Link>
                        </SidebarMenuButton>
                      </SidebarMenuItem>
                    );
                  })}
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          ))}
        </nav>
      </SidebarContent>

      <SidebarFooter className="px-3 pb-4">
        {!collapsed ? (
          <p className="text-[11px] leading-relaxed text-sidebar-foreground/45">
            RPA Performance &amp; Monitoring Platform
          </p>
        ) : null}
      </SidebarFooter>
    </Sidebar>
  );
}
