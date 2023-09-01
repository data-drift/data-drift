import { Params, useLoaderData } from "react-router-dom";
import {
  WaterfallChart,
  WaterfallChartProps,
} from "../components/Charts/WaterfallChart";
import { theme } from "../theme";
import {
  PeriodReport,
  Timegrain,
  TimegrainString,
  assertStringIsTimgrainString,
  getMetricReport,
  getTimegrainFromString,
} from "../services/data-drift";
import { getNiceTickValues } from "recharts-scale";

const getMetricCohortsData = async ({
  params,
}: {
  params: Params<string>;
}): Promise<WaterfallChartProps> => {
  const typedParams = assertParamsHasNeededProperties(params);
  const result = await getMetricReport(typedParams);
  const metricMetadata = result.data[typedParams.timegrainValue];
  const { data } = getWaterfallChartPropsFromMetadata(metricMetadata);
  return { data };
};

type Mutable<T> = {
  -readonly [P in keyof T]: T[P];
};

const getWaterfallChartPropsFromMetadata = (
  cohortMetric: PeriodReport
): WaterfallChartProps => {
  console.log(cohortMetric);
  // let latestValue = parseFloat(cohortMetric.InitialValue);

  const data = [] as Mutable<WaterfallChartProps["data"]>;

  const historyEntries = Object.entries(cohortMetric.History);
  historyEntries.sort(([_aCommitSha, aHistory], [_bCommitSha, bHistory]) => {
    return aHistory.CommitTimestamp - bHistory.CommitTimestamp;
  });
  historyEntries.forEach(([_commitSha, commit], index) => {
    const commitDate = new Date(commit.CommitTimestamp * 1000);
    const formatedDate = `${String(commitDate.getMonth() + 1).padStart(
      2,
      "0"
    )}-${String(commitDate.getDate()).padStart(2, "0")}`;
    if (index === 0) {
      const yMin = Math.min(
        ...historyEntries.map(([_sha, commit]) => parseFloat(commit.KPI))
      );
      const yMax = Math.max(
        ...historyEntries.map(([_sha, commit]) => parseFloat(commit.KPI))
      );
      const niceTicks = getNiceTickValues([yMin, yMax], 5);
      data.push({
        day: formatedDate,
        drift: [niceTicks[0], parseFloat(commit.KPI)],
        fill: theme.colors.text,
      });
    } else {
      const latestValue = parseFloat(historyEntries[index - 1][1].KPI);
      const newValue = parseFloat(commit.KPI);
      if (latestValue === newValue) {
        return;
      }
      data.push({
        day: formatedDate,
        drift: [latestValue, newValue],
        fill:
          latestValue > newValue ? theme.colors.dataDown : theme.colors.dataUp,
      });
    }
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
