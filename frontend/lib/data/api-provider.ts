import type { RpaDataProvider } from "@/lib/data/types";

function unavailable(): never {
  throw new Error("RPMP API data provider is not implemented in FE-05.");
}

export const apiDataProvider: RpaDataProvider = {
  async getUseCases() {
    return unavailable();
  },
  async getUseCase() {
    return unavailable();
  },
  async getSummary() {
    return unavailable();
  },
  async getPerformanceSeries() {
    return unavailable();
  },
  async getErrorDistribution() {
    return unavailable();
  },
  async getAutomationTypeSplit() {
    return unavailable();
  },
  async getReportHistory() {
    return unavailable();
  },
  async getEmailHistory() {
    return unavailable();
  },
  async getRecipients() {
    return unavailable();
  },
  async getReportTypes() {
    return unavailable();
  },
  async getReportPeriods() {
    return unavailable();
  },
};
