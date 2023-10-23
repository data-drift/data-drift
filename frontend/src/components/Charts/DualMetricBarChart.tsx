import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
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
    <ResponsiveContainer width="100%" height={300}>
      <BarChart
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
    </ResponsiveContainer>
  );
};

export default DualMetricBarChart;
