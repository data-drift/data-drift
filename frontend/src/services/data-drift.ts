import axios, { AxiosResponse } from "axios";
import { CommitParam } from "../pages/DisplayCommit/DisplayCommit";
import { MetricCohortsResults } from "./data-drift.types";
import { Endpoints } from "@octokit/types";

const DATA_DRIFT_API_URL =
  String(import.meta.env.VITE_DATADRIFT_SERVER_URL) ||
  "https://data-drift.herokuapp.com";

export const getPatchAndHeader = async (
  params: CommitParam & { installationId: string }
) => {
  const result = await axios.get<{
    patch: string;
    headers: string[];
    commitLink: string;
    date: string;
    filename: string;
    patchToLarge: boolean;
  }>(
    `${DATA_DRIFT_API_URL}/gh/${params.owner}/${params.repo}/commit/${params.commitSHA}`,
    { headers: { "Installation-Id": params.installationId } }
  );
  return {
    patch: result.data.patch,
    headers: result.data.headers,
    commitLink: result.data.commitLink,
    date: new Date(result.data.date),
    filename: result.data.filename,
    patchToLarge: result.data.patchToLarge,
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

const githubCommitListCache = new Map<
  string,
  AxiosResponse<
    Endpoints["GET /repos/{owner}/{repo}/commits"]["response"]["data"]
  >
>();

export const getCommitList = async (
  params: {
    installationId: string;
    owner: string;
    repo: string;
  },
  date?: string
) => {
  const cacheKey = `${params.owner}-${params.repo}-${params.installationId}-${
    date || "no-date"
  }`;

  if (githubCommitListCache.has(cacheKey)) {
    const cachedResult = githubCommitListCache.get(cacheKey);
    if (cachedResult) {
      return cachedResult;
    }
  }
  const result = await axios.get<
    Endpoints["GET /repos/{owner}/{repo}/commits"]["response"]["data"]
  >(`${DATA_DRIFT_API_URL}/gh/${params.owner}/${params.repo}/commits`, {
    headers: { "Installation-Id": params.installationId },
    params: { date },
  });

  githubCommitListCache.set(cacheKey, result);
  return result;
};

export const getCommitListLocalStrategy = async (
  tableName: string,
  date?: string
) => {
  const store = "default";

  const result = await axios.get<
    Endpoints["GET /repos/{owner}/{repo}/commits"]["response"]["data"]
  >(`${DATA_DRIFT_API_URL}/stores/${store}/tables/${tableName}/commits`, {
    params: { date },
  });

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

export type TimegrainAndDimensionString = `${TimegrainString}-${string}`;

export type MetricReport = Record<
  TimegrainString | TimegrainAndDimensionString,
  PeriodReport | undefined
>;

type CommitSha = string;
export interface PeriodReport {
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
  CommitDate: string;
  IsAfterPeriod: boolean;
  CommitUrl: string;
  CommitComments: CommitComment[] | null;
}

interface CommitComment {
  CommentAuthor: string;
  CommentBody: string;
}

export const ddCommitDiffUrlFactory = (params: {
  installationId: string;
  owner: string;
  repo: string;
  commitSha: string;
}) => {
  return `/report/${params.installationId}/${params.owner}/${params.repo}/commit/${params.commitSha}`;
};

export type DDConfigMetric = {
  filepath: string;
  upstreamFiles?: string[] | null;
  dateColumnName: string;
  KPIColumnName: string;
  metricName: string;
  timeGrains: string[];
  dimensions: string[];
};

export type DDConfig = {
  metrics: DDConfigMetric[];
};

const getConfigFromApi = async (params: {
  installationId: string;
  owner: string;
  repo: string;
}) => {
  const result = await axios.get<{ config: DDConfig }>(
    `${DATA_DRIFT_API_URL}/config/${params.owner}/${params.repo}`,
    { headers: { "Installation-Id": params.installationId } }
  );
  return result.data.config;
};

const configNameBuilder = (owner: string, repo: string) => {
  return `config-${owner}/${repo}`;
};

const getConfigFromSessionStorage = (
  owner: string,
  repo: string
): DDConfig | null => {
  const config = sessionStorage.getItem(configNameBuilder(owner, repo));
  if (config) {
    return JSON.parse(config) as DDConfig;
  } else {
    return null;
  }
};

const storeConfigInSessionStorage = (
  owner: string,
  repo: string,
  config: DDConfig
) => {
  sessionStorage.setItem(
    configNameBuilder(owner, repo),
    JSON.stringify(config)
  );
};

export const getConfig = async (params: {
  installationId: string;
  owner: string;
  repo: string;
}) => {
  const configFromStorage = getConfigFromSessionStorage(
    params.owner,
    params.repo
  );
  if (configFromStorage) {
    return configFromStorage;
  }
  const configFromApi = await getConfigFromApi(params);
  storeConfigInSessionStorage(params.owner, params.repo, configFromApi);
  return configFromApi;
};

export const getTableList = async () => {
  const result = await axios.get<{
    store: string;
    tables: string[];
  }>(`${DATA_DRIFT_API_URL}/stores/default/tables`);
  return result;
};

export const getTable = async (tableName: string) => {
  const result = await axios.get<{
    commits: {
      Message: string;
      Date: string;
      Sha: string;
    }[];
    store: string;
    table: string;
    tableColumns: string[];
  }>(`${DATA_DRIFT_API_URL}/stores/default/tables/${tableName}`);
  return result;
};

export const getMetricHistory = async (params: {
  store: string;
  table: string;
  metric: string;
  periodKey: string;
}) => {
  const { store, table, metric, periodKey } = params;
  const result = await axios.post<{
    metricHistory: {
      LineCount: number;
      Metric: string;
      IsMeasureAfterPeriod: boolean;
      MeasurementMetaData: {
        MeasurementTimestamp: number;
        MeasurementDate: string;
        MeasurementDateTime: string;
        MeasurementComments: {
          CommentAuthor: string;
          CommentBody: string;
        }[];

        MeasurementId: string;
      };
    }[];
    periodKey: string;
    store: string;
    table: string;
  }>(`${DATA_DRIFT_API_URL}/stores/${store}/tables/${table}/metrics`, {
    metric,
    period: periodKey,
  });
  return result;
};
