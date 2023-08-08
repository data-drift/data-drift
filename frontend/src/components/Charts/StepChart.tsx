import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
} from "recharts";

const formatXAxisTick = (tickValue: number) => {
  return `${Math.round(tickValue).toString()}d`;
};

const formatYAxisTick = (tickValue: number) => {
  return `${tickValue}%`;
};

export const StepChart = ({ data }: { data: any[] }) => {
  return (
    <LineChart
      width={500}
      height={300}
      data={data}
      margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
    >
      <CartesianGrid strokeDasharray="3 3" />
      <XAxis dataKey="x" tickFormatter={formatXAxisTick} type="number" />
      <YAxis tickFormatter={formatYAxisTick} />
      <Tooltip />
      <Legend />
      <Line
        type="stepAfter"
        dataKey="y"
        stroke="#8884d8"
        activeDot={{ r: 8 }}
      />
    </LineChart>
  );
};
