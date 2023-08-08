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

const formatXAxisTick = (tickValue: number) => {
  return `${Math.round(tickValue).toString()}d`;
};

const formatYAxisTick = (tickValue: number) => {
  return `${tickValue}%`;
};

export type YearMonthString = `${number}-${number}`;

type MetricEvolution = Array<
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
      width={500}
      height={300}
      data={data}
      margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
    >
      <CartesianGrid strokeDasharray="3 3" />
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
          activeDot={{ r: 8 }}
        />
      ))}
    </LineChart>
  );
};
