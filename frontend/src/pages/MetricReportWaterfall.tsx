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

const getMetricCohortsData = async ({
  params,
}: {
  params: Params<string>;
}): Promise<WaterfallChartProps> => {
  const typedParams = assertParamsHasNeededProperties(params);
  const result = await getMetricCohorts(typedParams);
  const metricMetadata =
    result.data.cohortsMetricsMetadata[typedParams.timegrainValue];
  console.log(metricMetadata);
  const data = [
    {
      day: "05-01",
      drift: [980, 1000],
      fill: theme.colors.text,
      isInitial: true,
    },
    {
      day: "05-02",
      drift: [1000, 1050],
      fill: theme.colors.dataUp,
    },
    {
      day: "05-03",
      drift: [1050, 1020],
      fill: theme.colors.dataDown,
    },
    {
      day: "05-04",
      drift: [1020, 1010],
      fill: theme.colors.dataDown,
    },
  ] as const;
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
