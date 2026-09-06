"use client";

import { useState } from "react";
import { LogOut } from "lucide-react";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";

export function LogoutControl() {
  const router = useRouter();
  const [status, setStatus] = useState<"idle" | "pending" | "error">("idle");

  async function logout() {
    if (status === "pending") {
      return;
    }

    setStatus("pending");

    try {
      const response = await fetch("/api/v1/auth/logout", {
        method: "POST",
        credentials: "include",
      });

      if (!response.ok) {
        setStatus("error");
        return;
      }

      router.replace("/login");
      router.refresh();
    } catch {
      setStatus("error");
    }
  }

  return (
    <div className="flex flex-col items-end gap-1">
      <Button
        aria-label="Log out"
        className="text-muted-foreground hover:bg-muted hover:text-foreground"
        disabled={status === "pending"}
        onClick={logout}
        size="icon"
        type="button"
        variant="ghost"
      >
        <LogOut aria-hidden="true" />
      </Button>
      {status === "error" ? (
        <p role="alert" className="text-xs text-destructive">
          Logout failed. Try again.
        </p>
      ) : null}
    </div>
  );
}
