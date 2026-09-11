import { ArrowLeft } from "lucide-react";
import Link from "next/link";

import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function UseCasesResourceState({
  detail = false,
  sourceUnavailable = false,
}: {
  detail?: boolean;
  sourceUnavailable?: boolean;
}) {
  return (
    <main className="mx-auto grid min-h-[50vh] w-full max-w-2xl place-items-center px-4 py-12 text-center">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">
          Use cases could not be loaded
        </h1>
        <p role="status" className="mt-2 text-sm text-muted-foreground">
          {sourceUnavailable
            ? "The stored RPMP data source is unavailable. Try again later."
            : "The use case data could not be loaded. Try again later."}
        </p>
        {detail ? (
          <Link
            className={cn(
              buttonVariants({ variant: "outline", size: "sm" }),
              "mt-5 rounded-full",
            )}
            href="/use-cases"
          >
            <ArrowLeft aria-hidden="true" className="mr-1 h-4 w-4" />
            Back to Use Cases
          </Link>
        ) : null}
      </div>
    </main>
  );
}
