import axios from "axios";
import { CommitInfo } from "../pages/DisplayCommit";
import { MetricCohortsResults } from "./data-drift.types";

const DATA_DRIFT_API_URL = "https://data-drift.herokuapp.com";

export const getPatchAndHeader = async (
  params: CommitInfo & { installationId: string }
) => {
  const result = await axios.get<{ patch: string; headers: string[] }>(
    `${DATA_DRIFT_API_URL}/gh/${params.owner}/${params.repo}/commit/${params.commitSHA}`,
    { headers: { "Installation-Id": params.installationId } }
  );
  return {
    patch: result.data.patch,
    headers: result.data.headers,
  };
};

export const getMetricCohorts = async ({
  installationId,
  metricName,
  timegrain,
}: {
  installationId: string;
  metricName: string;
  timegrain: Timegrain;
}) => {
  const result = await axios.get<MetricCohortsResults>(
    `${DATA_DRIFT_API_URL}/metrics/${metricName}/cohorts/${timegrain}`,
    { headers: { "Installation-Id": installationId } }
  );
  return result;
};

// Define the custom type
export type Timegrain = "year" | "quarter" | "month" | "week" | "day";

// The assertion function
export function assertTimegrain(value: string): asserts value is Timegrain {
  if (
    value !== "year" &&
    value !== "quarter" &&
    value !== "month" &&
    value !== "week" &&
    value !== "day"
  ) {
    throw new Error("Value is not a valid time unit!");
  }
}
