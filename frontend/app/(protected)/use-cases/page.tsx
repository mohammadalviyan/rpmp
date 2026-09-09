import type { Metadata } from "next";

import { UseCasesList } from "@/components/use-cases-list";
import { getRpaDataProvider } from "@/lib/data/provider";

export const metadata: Metadata = {
  title: "Use Cases | RPMP",
  description:
    "Browse monitored use cases by automation type, status, success rate, and issues.",
};

export default async function UseCasesPage() {
  const useCases = await getRpaDataProvider().getUseCases();

  return <UseCasesList useCases={useCases} />;
}
