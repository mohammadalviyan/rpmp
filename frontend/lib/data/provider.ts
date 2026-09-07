import { apiDataProvider } from "@/lib/data/api-provider";
import { mockDataProvider } from "@/lib/data/mock-provider";
import type { RpaDataProvider } from "@/lib/data/types";

export type RpaDataMode = "api" | "mock";

export function getRpaDataMode(
  value = process.env.RPMP_DATA_MODE,
): RpaDataMode {
  return value === "api" ? "api" : "mock";
}

export function getRpaDataProvider(
  mode = getRpaDataMode(),
): RpaDataProvider {
  return mode === "api" ? apiDataProvider : mockDataProvider;
}
