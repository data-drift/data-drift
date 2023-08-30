import { Params, useLoaderData } from "react-router-dom";
import {
  WaterfallChart,
  WaterfallChartProps,
} from "../components/Charts/WaterfallChart";
import { theme } from "../theme";
import {
  Timegrain,
  TimegrainString,
  assertStringIsTimgrainString,
  getMetricCohorts,
  getTimegrainFromString,
} from "../services/data-drift";
import { CohortMetric } from "../services/data-drift.types";

const getMetricCohortsData = async ({
  params,
}: {
  params: Params<string>;
}): Promise<WaterfallChartProps> => {
  const typedParams = assertParamsHasNeededProperties(params);
  const result = await getMetricCohorts(typedParams);
  const metricMetadata =
    result.data.cohortsMetricsMetadata[typedParams.timegrainValue];
  const { data } = getWaterfallChartPropsFromMetadata(metricMetadata);
  return { data };
};

type Mutable<T> = {
  -readonly [P in keyof T]: T[P];
};

const getWaterfallChartPropsFromMetadata = (
  cohortMetric: CohortMetric
): WaterfallChartProps => {
  console.log(cohortMetric);
  let latestValue = parseFloat(cohortMetric.InitialValue);

  const data = [] as Mutable<WaterfallChartProps["data"]>;
  Object.keys(cohortMetric.RelativeHistory).forEach((cohortKey) => {
    const aberanteValue =
      parseFloat(cohortMetric.RelativeHistory[cohortKey].RelativeValue) > 100;
    if (aberanteValue) {
      return;
    }
    const newValue =
      parseFloat(cohortMetric.RelativeHistory[cohortKey].RelativeValue) *
      parseFloat(cohortMetric.InitialValue);

    if (newValue == latestValue) {
      return;
    }
    const result = {
      day: "05-04",
      drift: [latestValue, newValue],
      fill:
        latestValue > newValue ? theme.colors.dataDown : theme.colors.dataUp,
    } as const;
    latestValue = newValue;
    data.push(result);
  });
  return { data };
};

function assertParamsHasNeededProperties(params: Params<string>): {
  installationId: string;
  metricName: string;
  timegrain: Timegrain;
  timegrainValue: TimegrainString;
} {
  const { installationId, metricName, timegrainValue } = params;
  if (!installationId || !metricName || !timegrainValue) {
    throw new Error("Invalid params");
  }
  assertStringIsTimgrainString(timegrainValue);
  const timegrain = getTimegrainFromString(timegrainValue);

  return { installationId, metricName, timegrain, timegrainValue };
}

const MetricReportWaterfall = () => {
  const { data } = useLoaderData() as WaterfallChartProps;

  return <WaterfallChart data={data} />;
};

MetricReportWaterfall.loader = getMetricCohortsData;

export default MetricReportWaterfall;
