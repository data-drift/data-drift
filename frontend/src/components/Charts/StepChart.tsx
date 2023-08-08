import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
} from "recharts";

import { scaleLinear } from "d3-scale";

const formatXAxisTick = (tickValue: number) => {
  return `${Math.round(tickValue).toString()}d`;
};

const formatYAxisTick = (tickValue: number) => {
  return `${tickValue}%`;
};

const colorSelector = (year: string) => {
  switch (year) {
    case "2022":
      return ["red", "blue"];
    case "2023":
      return ["green", "yellow"];
    default:
      return ["black", "white"];
  }
};

type YearMonthString = `${number}-${number}`;

type MetricEvolution = Array<
  {
    daysSinceFirstReport: number;
  } & Record<YearMonthString, number>
>;

const getMetricColor = (yearMonthString: YearMonthString) => {
  const [year, month] = yearMonthString.split("-");
  console.log(year, month);
  const scale = scaleLinear([0, 11], colorSelector(year));

  console.log(scale(parseInt(month, 10)));
  return scale(parseInt(month, 10));
};

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
