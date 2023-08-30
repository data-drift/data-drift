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

type YearString = `${number}`;
type YearMonthString = `${number}-${string & { length: 2 }}`;
type YearMonthDayString = `${number}-${string & { length: 2 }}-${string & {
  length: 2;
}}`;
type YearWeekString = `${number}-W${
  | (number & { length: 1 })
  | (string & { length: 2 })}`;
type YearQuarterString = `${number}-Q${1 | 2 | 3 | 4}`;
export type TimegrainString =
  | YearString
  | YearMonthString
  | YearMonthDayString
  | YearWeekString
  | YearQuarterString;

export function assertStringIsTimgrainString(
  str: string
): asserts str is TimegrainString {
  if (
    str.match(/^\d{4}$/) !== null ||
    str.match(/^\d{4}-\d{2}$/) !== null ||
    str.match(/^\d{4}-\d{2}-\d{2}$/) !== null ||
    str.match(/^\d{4}-W\d{1,2}$/) !== null ||
    str.match(/^\d{4}-Q[1-4]$/) !== null
  ) {
    return;
  } else {
    throw new Error("Invalid timegrain string!");
  }
}

export function getTimegrainFromString(str: TimegrainString): Timegrain {
  if (str.match(/^\d{4}$/) !== null) {
    return "year";
  } else if (str.match(/^\d{4}-\d{2}$/) !== null) {
    return "month";
  } else if (str.match(/^\d{4}-\d{2}-\d{2}$/) !== null) {
    return "day";
  } else if (str.match(/^\d{4}-W\d{1,2}$/) !== null) {
    return "week";
  } else if (str.match(/^\d{4}-Q[1-4]$/) !== null) {
    return "quarter";
  } else {
    throw new Error("Invalid timegrain string!");
  }
}

export const getMetricReport = async ({
  installationId,
  metricName,
}: {
  installationId: string;
  metricName: string;
  timegrain: Timegrain;
}) => {
  const result = await axios.get<MetricReport>(
    `${DATA_DRIFT_API_URL}/metrics/${metricName}/reports`,
    { headers: { "Installation-Id": installationId } }
  );
  return result;
};

export type MetricReport = Record<TimegrainString, PeriodReport>;

type CommitSha = string;
interface PeriodReport {
  TimeGrain: Timegrain;
  Period: TimegrainString;
  Dimension: string;
  DimensionValue: string;
  History: { [key: CommitSha]: History };
}

interface History {
  Lines: number;
  KPI: string;
  CommitTimestamp: number;
  CommitUrl: string;
  CommitComments: CommitComment[] | null;
}

interface CommitComment {
  CommentAuthor: string;
  CommentBody: string;
}
