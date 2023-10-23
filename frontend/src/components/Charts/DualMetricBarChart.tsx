import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
} from "recharts";
import { theme } from "../../theme";

type Props = {
  data: Array<{
    name: string;
    before: number;
    after: number;
  }>;
};

const DualMetricBarChart = ({ data }: Props) => {
  return (
    <BarChart
      width={500}
      height={300}
      data={data}
      margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
    >
      <CartesianGrid strokeDasharray="3 3" />
      <XAxis dataKey="name" />
      <YAxis />
      <Tooltip />
      <Legend />
      <Bar dataKey="before" fill={theme.colors.strongNegative} />
      <Bar dataKey="after" fill={theme.colors.strongPositive} />
    </BarChart>
  );
};

export default DualMetricBarChart;
