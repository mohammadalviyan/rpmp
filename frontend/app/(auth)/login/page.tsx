import { redirect } from "next/navigation";

import { LoginForm } from "@/components/login-form";
import { getCurrentUser } from "@/lib/api/server";

export default async function LoginPage() {
  const user = await getCurrentUser();

  if (user) {
    redirect("/dashboard");
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-muted px-4 py-12">
      <section
        aria-labelledby="login-title"
        className="w-full max-w-md rounded-xl border bg-card p-8 text-card-foreground shadow-sm"
      >
        <div className="mb-8 space-y-2">
          <p className="text-sm font-semibold tracking-wide text-muted-foreground">
            RPMP
          </p>
          <h1 id="login-title" className="text-2xl font-semibold">
            Sign in
          </h1>
          <p className="text-sm text-muted-foreground">
            Use your Employee ID and password to continue.
          </p>
        </div>
        <LoginForm />
      </section>
    </main>
  );
}
