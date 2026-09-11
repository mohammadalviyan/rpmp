import type { Metadata } from "next";
import { redirect } from "next/navigation";

import { ApiUseCaseDetail } from "@/components/api-use-case-detail";
import {
  UseCaseDetail,
  UseCaseNotFound,
} from "@/components/use-case-detail";
import { UseCasesResourceState } from "@/components/use-cases-resource-state";
import { getUseCase } from "@/lib/api/server";
import {
  getRpaDataMode,
  getRpaDataProvider,
} from "@/lib/data/provider";

export const metadata: Metadata = {
  title: "Use Case Details | RPMP",
  description: "Review one use case and its stored performance.",
};

export default async function UseCaseDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  if (getRpaDataMode() === "api") {
    const result = await getUseCase(id);

    if (result.status === "unauthenticated") {
      redirect("/login");
    }

    if (result.status === "not_found") {
      return <UseCaseNotFound />;
    }

    if (result.status !== "success") {
      return (
        <UseCasesResourceState
          detail
          sourceUnavailable={result.status === "source_unavailable"}
        />
      );
    }

    return <ApiUseCaseDetail data={result.data} />;
  }

  const useCase = await getRpaDataProvider().getUseCase(id);

  return useCase ? <UseCaseDetail useCase={useCase} /> : <UseCaseNotFound />;
}
