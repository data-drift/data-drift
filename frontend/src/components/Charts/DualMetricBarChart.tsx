import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  Rectangle,
  RectangleProps,
} from "recharts";
import { theme } from "../../theme";

type Props = {
  data: Array<{
    name: string;
    before: number;
    after: number;
  }>;
};

const CustomBarShape = (props: RectangleProps) => {
  const { fill, x, y, width } = props;

  const brightFill = fill
    ? fill.slice(0, fill.lastIndexOf(",")) + ",1)"
    : "white";

  return (
    <>
      <Rectangle {...props} y={y} fill={fill} />
      <Rectangle x={x} y={y} width={width} height={2} fill={brightFill} />
    </>
  );
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
        <Bar
          dataKey="before"
          fill={theme.colors.hexToRgba(theme.colors.strongNegative, 0.4)}
          shape={<CustomBarShape />}
        />
        <Bar
          dataKey="after"
          fill={theme.colors.hexToRgba(theme.colors.strongPositive, 0.4)}
          shape={<CustomBarShape />}
        />
      </BarChart>
    </ResponsiveContainer>
  );
};

export default DualMetricBarChart;
