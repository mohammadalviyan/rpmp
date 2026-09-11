import type { Metadata } from "next";
import { redirect } from "next/navigation";

import { ApiUseCasesList } from "@/components/api-use-cases-list";
import { UseCasesList } from "@/components/use-cases-list";
import { UseCasesResourceState } from "@/components/use-cases-resource-state";
import { getUseCases } from "@/lib/api/server";
import {
  getRpaDataMode,
  getRpaDataProvider,
} from "@/lib/data/provider";

export const metadata: Metadata = {
  title: "Use Cases | RPMP",
  description: "Browse monitored use cases and their stored performance.",
};

export default async function UseCasesPage() {
  if (getRpaDataMode() === "api") {
    const result = await getUseCases();

    if (result.status === "unauthenticated") {
      redirect("/login");
    }

    if (result.status !== "success") {
      return (
        <UseCasesResourceState
          sourceUnavailable={result.status === "source_unavailable"}
        />
      );
    }

    return <ApiUseCasesList useCases={result.data.items} />;
  }

  const useCases = await getRpaDataProvider().getUseCases();

  return <UseCasesList useCases={useCases} />;
}
