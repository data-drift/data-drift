import {
  StepChartProps,
  YearMonthString,
} from "../components/Charts/StepChart";
import { CohortPoint, CohortsMetricsMetadata } from "./data-drift.types";

export const mapCohortsMetricsMetadataToStepChartProps = (
  cohortMetricMetadata: CohortsMetricsMetadata
): StepChartProps => {
  const metricNames: StepChartProps["metricNames"] = [];
  const data: StepChartProps["data"] = [];
  for (const metricName of Object.keys(
    cohortMetricMetadata
  ) as YearMonthString[]) {
    const cohort = cohortMetricMetadata[metricName];
    const RelativeHistory = sortRelativeHistory(cohort.RelativeHistory);

    metricNames.push(metricName);
    let latestKPI = undefined as number | undefined;
    RelativeHistory.sort();

    data.push({
      daysSinceFirstReport: 0,
      [metricName]: 0,
    });

    RelativeHistory.forEach((cohortMetric) => {
      const currentKPI = parseFloat(cohortMetric.RelativeValue);
      if (latestKPI === undefined || currentKPI != latestKPI) {
        const datum = {
          [metricName]: currentKPI,
          daysSinceFirstReport: parseFloat(cohortMetric.DaysFromHistorization),
        };
        data.push(datum);
        latestKPI = currentKPI;
      }
    });
  }

  return { data, metricNames };
};

function sortRelativeHistory(
  cohort: Record<string, CohortPoint>
): CohortPoint[] {
  const relativeHistoryArray = Object.keys(cohort).map((key) => ({
    ...cohort[key],
    daysFromHistorisation: parseFloat(cohort[key].DaysFromHistorization),
    relativeValue: parseFloat(cohort[key].RelativeValue),
  }));

  relativeHistoryArray.sort(
    (a, b) => a.daysFromHistorisation - b.daysFromHistorisation
  );

  return relativeHistoryArray;
}
