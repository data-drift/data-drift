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
    // if (metricName !== "2023-02") continue;
    const cohort = cohortMetricMetadata[metricName];
    const RelativeHistory = sortRelativeHistory(cohort.RelativeHistory);
    console.log(cohort.RelativeHistory);
    if (RelativeHistory.length < 2) {
      continue;
    }
    metricNames.push(metricName);
    let latestKPI = undefined as number | undefined;
    console.log(RelativeHistory);
    RelativeHistory.sort();
    console.log(RelativeHistory);

    RelativeHistory.forEach((cohortMetric) => {
      const currentKPI = parseFloat(cohortMetric.RelativeValue);
      if (latestKPI === undefined || currentKPI != latestKPI) {
        console.log("days", parseFloat(cohortMetric.DaysFromHistorization));
        const datum = {
          [metricName]: currentKPI,
          daysSinceFirstReport: parseFloat(cohortMetric.DaysFromHistorization),
        };
        data.push(datum);
        latestKPI = currentKPI;
      }
    });
  }
  console.log(data);

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
