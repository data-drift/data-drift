import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
} from "recharts";
import { getMetricColor } from "./colors.utils";
import { theme } from "../../theme";

const formatXAxisTick = (tickValue: number) => {
  return `${Math.round(tickValue).toString()}d`;
};

const formatYAxisTick = (tickValue: number) => {
  return `${tickValue}%`;
};

export type YearMonthString = `${number}-${number}`;

export type MetricEvolution = Array<
  {
    daysSinceFirstReport: number;
  } & Record<YearMonthString, number>
>;

export const StepChart = ({
  data,
  metricNames,
}: {
  data: MetricEvolution;
  metricNames: YearMonthString[];
}) => {
  return (
    <LineChart
      width={1000}
      height={600}
      data={data}
      margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
    >
      <CartesianGrid strokeDasharray="3 3" fill={theme.colors.background2} />
      <XAxis
        dataKey="daysSinceFirstReport"
        tickFormatter={formatXAxisTick}
        type="number"
      />
      <YAxis tickFormatter={formatYAxisTick} />
      <Tooltip />
      <Legend />
      {metricNames.map((metricName) => (
        <Line
          type="stepAfter"
          dataKey={metricName}
          stroke={getMetricColor(metricName)}
          dot={false}
        />
      ))}
    </LineChart>
  );
};
