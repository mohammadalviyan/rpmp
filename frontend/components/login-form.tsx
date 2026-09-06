"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import type { ApiError, LoginRequest } from "@/lib/api/types";

const INVALID_CREDENTIALS =
  "Employee ID or password is incorrect.";
const UNAVAILABLE =
  "RPMP is unavailable right now. Please try again in a few minutes.";
const UNEXPECTED =
  "We could not sign you in. Check your details and try again.";

type FormStatus = "idle" | "pending" | "invalid" | "unavailable" | "unexpected";

function statusMessage(status: FormStatus) {
  switch (status) {
    case "invalid":
      return INVALID_CREDENTIALS;
    case "unavailable":
      return UNAVAILABLE;
    case "unexpected":
      return UNEXPECTED;
    default:
      return null;
  }
}

export function LoginForm() {
  const router = useRouter();
  const [status, setStatus] = useState<FormStatus>("idle");
  const pending = status === "pending";
  const message = statusMessage(status);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (pending) {
      return;
    }

    setStatus("pending");
    const form = new FormData(event.currentTarget);
    const payload: LoginRequest = {
      employee_id: String(form.get("employee_id") ?? ""),
      password: String(form.get("password") ?? ""),
    };

    try {
      const response = await fetch("/api/v1/auth/login", {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        router.replace("/dashboard");
        router.refresh();
        return;
      }

      let error: Partial<ApiError> = {};
      try {
        error = (await response.json()) as Partial<ApiError>;
      } catch {
        // A non-JSON response is handled as an unavailable or unexpected error.
      }

      if (response.status >= 500 || error.code === "internal_error" || error.code === "source_unavailable") {
        setStatus("unavailable");
      } else if (error.code === "invalid_credentials") {
        setStatus("invalid");
      } else {
        setStatus("unexpected");
      }
    } catch {
      setStatus("unavailable");
    }
  }

  return (
    <form className="space-y-5" onSubmit={handleSubmit}>
      <div className="space-y-2">
        <Label htmlFor="employee_id">Employee ID</Label>
        <Input
          className="h-10"
          id="employee_id"
          name="employee_id"
          type="text"
          autoComplete="username"
          inputMode="numeric"
          required
          disabled={pending}
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="password">Password</Label>
        <Input
          className="h-10"
          id="password"
          name="password"
          type="password"
          autoComplete="current-password"
          required
          disabled={pending}
        />
      </div>

      {message ? (
        <p role="alert" className="text-sm text-destructive">
          {message}
        </p>
      ) : null}

      {/* Login-only hover so the CTA lands on Cakrawala 600 instead of primary/80. */}
      <Button
        className="h-10 w-full text-sm font-semibold hover:bg-[var(--bri-primary-cakrawala-600)]"
        type="submit"
        disabled={pending}
      >
        {pending ? "Signing in..." : "Sign in"}
      </Button>
    </form>
  );
}
