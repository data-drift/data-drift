import { useLoaderData } from "react-router-dom";
import {
  WaterfallChart,
  WaterfallChartProps,
} from "../components/Charts/WaterfallChart";
import { theme } from "../theme";

const getMetricCohortsData = (): WaterfallChartProps => {
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

const MetricReportWaterfall = () => {
  const { data } = useLoaderData() as WaterfallChartProps;

  return <WaterfallChart data={data} />;
};

MetricReportWaterfall.loader = getMetricCohortsData;

export default MetricReportWaterfall;
