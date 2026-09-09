import type { Metadata } from "next";

import {
  UseCaseDetail,
  UseCaseNotFound,
} from "@/components/use-case-detail";
import { getRpaDataProvider } from "@/lib/data/provider";

export const metadata: Metadata = {
  title: "Use Case Details | RPMP",
  description: "Review one use case's performance and issue history.",
};

export default async function UseCaseDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  const useCase = await getRpaDataProvider().getUseCase(id);

  return useCase ? <UseCaseDetail useCase={useCase} /> : <UseCaseNotFound />;
}
