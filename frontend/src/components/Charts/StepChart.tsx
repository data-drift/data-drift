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
import { useState } from "react";

const formatXAxisTick = (tickValue: number) => {
  return `${Math.round(tickValue).toString()}d`;
};

const formatYAxisTick = (tickValue: number) => {
  return `${tickValue}%`;
};

const formatToolTipValueTick = (tickValue: number) => {
  return `${tickValue.toPrecision(3)}%`;
};

export type YearMonthString = `${number}-${number}`;

export type MetricEvolution = Array<
  {
    daysSinceFirstReport: number;
  } & Record<YearMonthString, number>
>;

export type StepChartProps = {
  data: MetricEvolution;
  metricNames: YearMonthString[];
};

export const StepChart = ({ data, metricNames }: StepChartProps) => {
  const [highlightedMetric, setHighlightedMetric] = useState<string | null>(
    null
  );

  return (
    <LineChart
      width={1000}
      height={600}
      data={data}
      margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
    >
      <CartesianGrid
        fill={theme.colors.background2}
        vertical={false}
        strokeDasharray={"2 2"}
      />
      <XAxis
        dataKey="daysSinceFirstReport"
        tickFormatter={formatXAxisTick}
        type="number"
      />
      <YAxis tickFormatter={formatYAxisTick} />
      <Tooltip
        formatter={formatToolTipValueTick}
        labelFormatter={formatXAxisTick}
        contentStyle={{ backgroundColor: theme.colors.background }}
      />
      <Legend
        onMouseEnter={({ dataKey }) => setHighlightedMetric(dataKey as string)}
        onMouseLeave={() => setHighlightedMetric(null)}
      />
      {metricNames.map((metricName) => (
        <Line
          type="stepAfter"
          dataKey={metricName}
          stroke={getMetricColor(metricName)}
          strokeWidth={highlightedMetric === metricName ? 5 : 1}
          dot={false}
        />
      ))}
    </LineChart>
  );
};
