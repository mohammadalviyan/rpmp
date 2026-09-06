"use client";

import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import {
  Boxes,
  FileChartColumn,
  History,
  LayoutDashboard,
  Mail,
  MailCheck,
  Menu,
  ServerCog,
  SlidersHorizontal,
  Users,
  X,
  type LucideIcon,
} from "lucide-react";
import Link from "next/link";

import { cn } from "@/lib/utils";

type NavItem = {
  label: string;
  icon: LucideIcon;
  href?: string;
};

type NavSection = {
  title: string;
  items: NavItem[];
};

// Only Overview has a route. The rest are placeholders for menus RPMP has not
// built yet, so they stay disabled instead of linking anywhere.
const navSections: NavSection[] = [
  {
    title: "Dashboard",
    items: [{ label: "Overview", icon: LayoutDashboard, href: "/dashboard" }],
  },
  {
    title: "Report Management",
    items: [
      { label: "Use Cases", icon: Boxes },
      { label: "Generate Report", icon: FileChartColumn },
      { label: "Report History", icon: History },
    ],
  },
  {
    title: "Email",
    items: [
      { label: "Email Management", icon: Mail },
      { label: "Email History", icon: MailCheck },
    ],
  },
  {
    title: "Settings",
    items: [
      { label: "General Settings", icon: SlidersHorizontal },
      { label: "User Management", icon: Users },
      { label: "Email Configuration", icon: ServerCog },
    ],
  },
];

type SidebarContextValue = {
  open: boolean;
  setOpen: (open: boolean) => void;
};

const SidebarContext = createContext<SidebarContextValue | null>(null);

function useSidebar(): SidebarContextValue {
  const context = useContext(SidebarContext);

  if (!context) {
    throw new Error("Sidebar components must render inside SidebarProvider.");
  }

  return context;
}

export function SidebarProvider({ children }: { children: ReactNode }) {
  const [open, setOpen] = useState(false);

  return (
    <SidebarContext.Provider value={{ open, setOpen }}>
      {children}
    </SidebarContext.Provider>
  );
}

export function SidebarTrigger() {
  const { open, setOpen } = useSidebar();

  return (
    <button
      aria-controls="app-sidebar"
      aria-expanded={open}
      aria-label={open ? "Close menu" : "Open menu"}
      className="grid size-10 shrink-0 place-items-center rounded-xl border text-muted-foreground hover:bg-muted hover:text-foreground lg:hidden"
      onClick={() => setOpen(!open)}
      type="button"
    >
      {open ? (
        <X aria-hidden="true" className="size-5" />
      ) : (
        <Menu aria-hidden="true" className="size-5" />
      )}
    </button>
  );
}

export function AppSidebar() {
  const { open, setOpen } = useSidebar();

  useEffect(() => {
    if (!open) {
      return;
    }

    function closeOnEscape(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setOpen(false);
      }
    }

    document.addEventListener("keydown", closeOnEscape);
    return () => document.removeEventListener("keydown", closeOnEscape);
  }, [open, setOpen]);

  return (
    <>
      {open ? (
        <button
          aria-label="Close menu"
          className="fixed inset-0 z-30 bg-[var(--bri-black-opacity-50)] lg:hidden"
          onClick={() => setOpen(false)}
          type="button"
        />
      ) : null}

      <aside
        className={cn(
          "fixed inset-y-0 left-0 z-40 flex w-72 flex-col border-r bg-card transition-transform duration-200 ease-out lg:w-64 lg:translate-x-0",
          open ? "translate-x-0" : "-translate-x-full",
        )}
        id="app-sidebar"
      >
        <div className="flex items-center justify-between gap-3 px-5 py-5 lg:px-6 lg:py-7">
          <Link
            aria-label="RPMP Overview"
            className="flex items-center gap-3"
            href="/dashboard"
            onClick={() => setOpen(false)}
          >
            <span className="grid h-9 w-9 place-items-center rounded-xl bg-primary text-sm font-bold text-primary-foreground">
              R
            </span>
            <span>
              <span className="block text-base font-bold tracking-wide">RPMP</span>
              <span className="block text-[11px] text-muted-foreground">
                Performance monitoring
              </span>
            </span>
          </Link>
        </div>

        <nav
          aria-label="Primary navigation"
          className="flex-1 overflow-y-auto px-4 pb-6"
        >
          {navSections.map((section) => (
            <div key={section.title} className="mb-5 last:mb-0">
              <p className="px-4 pb-2 text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">
                {section.title}
              </p>
              <ul className="space-y-1">
                {section.items.map((item) => (
                  <li key={item.label}>
                    {item.href ? (
                      <Link
                        aria-current="page"
                        className="flex min-h-11 items-center gap-3 rounded-xl bg-primary px-4 text-sm font-semibold text-primary-foreground"
                        href={item.href}
                        onClick={() => setOpen(false)}
                      >
                        <item.icon aria-hidden="true" className="size-4" />
                        {item.label}
                      </Link>
                    ) : (
                      <button
                        aria-disabled="true"
                        className="flex min-h-11 w-full cursor-not-allowed items-center gap-3 rounded-xl px-4 text-sm font-medium text-muted-foreground opacity-60"
                        disabled
                        type="button"
                      >
                        <item.icon aria-hidden="true" className="size-4" />
                        {item.label}
                      </button>
                    )}
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </nav>
      </aside>
    </>
  );
}
