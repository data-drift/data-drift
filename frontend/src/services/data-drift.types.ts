type CohortDate = string;
type CohortDates = CohortDate[];
type TimestampString = string;

export interface MetricCohortsResults {
  cohortDates: CohortDates;
  cohortsMetricsMetadata: CohortsMetricsMetadata;
  dataIndexedByTimestamp: { [key: TimestampString]: DataIndexedByTimestamp };
  timegrain: string;
}

export type CohortsMetricsMetadata = Record<CohortDate, The202>;

export interface CohortPoint {
  RelativeValue: string;
  DaysFromHistorization: string;
  ComputationTimetamp: number;
}

export interface The202 {
  TimeGrain: string;
  PeriodKey: string;
  InitialValue: string;
  FirstDate: Date;
  RelativeHistory: { [key: string]: CohortPoint };
}

export type DataIndexedByTimestamp = Record<CohortDate, string>;
