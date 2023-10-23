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

type DualBarPayload = {
  name: string;
  before: number;
  after: number;
};
type Props = {
  data: Array<DualBarPayload>;
};

const CustomBarShape = (
  props: RectangleProps & { dataKey: "before" | "after" }
) => {
  const { fill, x, y, width, dataKey } = props;

  // Payload is present in props, in that scenario we can use it to get the value
  // @ts-ignore
  const payload = props.payload as DualBarPayload | undefined;
  console.log(dataKey);
  const { before, after } = payload ? payload : { before: null, after: null };
  const value =
    dataKey === "before" ? before : dataKey === "after" ? after : null;

  const brightFill = fill
    ? fill.slice(0, fill.lastIndexOf(",")) + ",1)"
    : "white";

  return (
    <>
      <Rectangle {...props} y={y} fill={fill} />
      <Rectangle x={x} y={y} width={width} height={2} fill={brightFill} />
      {value !== undefined &&
        value !== null &&
        x !== undefined &&
        width !== undefined &&
        y !== undefined && (
          <text
            x={x + width / 2}
            y={y + 12}
            fill={theme.colors.hexToRgba(theme.colors.text, 1)}
            textAnchor="middle"
            dominantBaseline="middle"
            fontSize={10}
          >
            {value.toLocaleString()}
          </text>
        )}
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
          shape={<CustomBarShape dataKey="before" />}
        />
        <Bar
          dataKey="after"
          fill={theme.colors.hexToRgba(theme.colors.strongPositive, 0.4)}
          shape={<CustomBarShape dataKey="after" />}
        />
      </BarChart>
    </ResponsiveContainer>
  );
};

export default DualMetricBarChart;
