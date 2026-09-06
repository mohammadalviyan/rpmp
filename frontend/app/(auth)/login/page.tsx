import Image from "next/image";
import { redirect } from "next/navigation";

import { LoginForm } from "@/components/login-form";
import { getCurrentUser } from "@/lib/api/server";

export default async function LoginPage() {
  const user = await getCurrentUser();

  if (user) {
    redirect("/dashboard");
  }

  return (
    <main className="grid min-h-screen grid-cols-1 bg-background lg:grid-cols-[1.15fr_1fr]">
      <div className="relative hidden lg:block">
        <Image
          src="/BRI-Login.webp"
          alt=""
          fill
          priority
          sizes="(min-width: 1024px) 55vw, 0px"
          className="object-cover"
        />
      </div>

      <div className="flex items-center justify-center px-4 py-12">
        <section
          aria-labelledby="login-title"
          className="w-full max-w-md overflow-hidden rounded-xl border bg-card text-card-foreground shadow-sm"
        >
          <div aria-hidden="true" className="h-1 bg-primary" />

          <div className="p-8">
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
          </div>
        </section>
      </div>
    </main>
  );
}
