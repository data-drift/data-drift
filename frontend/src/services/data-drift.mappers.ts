import {
  StepChartProps,
  YearMonthString,
} from "../components/Charts/StepChart";
import { CohortsMetricsMetadata } from "./data-drift.types";

export const mapCohortsMetricsMetadataToStepChartProps = (
  cohortMetricMetadata: CohortsMetricsMetadata
): StepChartProps => {
  const metricNames: StepChartProps["metricNames"] = Object.keys(
    cohortMetricMetadata
  ) as YearMonthString[];
  const data: StepChartProps["data"] = [];
  for (const metricName of metricNames) {
    const cohort = cohortMetricMetadata[metricName];
    const RelativeHistory = Object.keys(cohort.RelativeHistory);
    if (RelativeHistory.length < 2) {
      continue;
    }
    RelativeHistory.forEach((timeFromHistorisationInMs) => {
      const cohortMetric = cohort.RelativeHistory[timeFromHistorisationInMs];
      const datum = {
        [metricName]: parseFloat(cohortMetric.RelativeValue),
        daysSinceFirstReport: parseFloat(cohortMetric.DaysFromHistorization),
      };
      data.push(datum);
    });
  }

  return { data, metricNames };
};
