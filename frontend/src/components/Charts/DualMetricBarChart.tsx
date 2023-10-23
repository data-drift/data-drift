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
import TrendChip from "./TrendChip";

type DualBarPayload = {
  name: string;
  before: number;
  after: number;
  percentageChange: number;
};
type Props = {
  data: Array<DualBarPayload>;
};

function getShouldDisplayTrendChip(
  trend: "up" | "down" | "neutral",
  dataKey: "before" | "after"
): boolean {
  switch (trend) {
    case "up":
      return dataKey === "after";
    case "down":
      return dataKey === "before";
    case "neutral":
      return dataKey === "before";
  }
}
const CustomBarShape = (
  props: RectangleProps & { dataKey: "before" | "after" }
) => {
  const { fill, x, y, width, dataKey } = props;

  // Payload is present in props, in that scenario we can use it to get the value
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  const payload = props.payload as DualBarPayload | undefined;
  const { before, after, percentageChange } = payload
    ? payload
    : { before: 0, after: 0, percentageChange: 0 };
  const value =
    dataKey === "before" ? before : dataKey === "after" ? after : null;

  if (
    value == undefined ||
    value == null ||
    x == undefined ||
    width == undefined ||
    y == undefined
  ) {
    return null;
  }

  const trend = before < after ? "up" : before > after ? "down" : "neutral";
  const shouldDisplayTrendChip = getShouldDisplayTrendChip(trend, dataKey);
  const trendChipOffset = dataKey === "before" ? width : 0;

  const brightFill = fill
    ? fill.slice(0, fill.lastIndexOf(",")) + ",1)"
    : "white";

  return (
    <>
      <Rectangle {...props} y={y} fill={fill} />
      <Rectangle x={x} y={y} width={width} height={2} fill={brightFill} />

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
      {shouldDisplayTrendChip && (
        <foreignObject
          x={x - 100 + trendChipOffset}
          y={y - 40}
          width={200}
          height={50}
        >
          <div
            style={{
              display: "flex",
              justifyContent: "center",
            }}
          >
            <TrendChip
              trend={trend}
              absoluteValue={Math.abs(percentageChange)}
            />
          </div>
        </foreignObject>
      )}
    </>
  );
};

const barSize = 80;

const DualMetricBarChart = ({ data }: Props) => {
  const totalWidth = (barSize + 20) * data.length * 2 + 30; // Assuming 30px additional margin

  return (
    <ResponsiveContainer width={totalWidth} height={400}>
      <BarChart
        data={data}
        margin={{ top: 48, right: 30, left: 20, bottom: 5 }}
      >
        <CartesianGrid
          stroke={theme.colors.hexToRgba(theme.colors.background2, 0.8)}
          vertical={false}
        />
        <XAxis dataKey="name" />
        <YAxis />
        <Tooltip contentStyle={{ backgroundColor: theme.colors.background }} />
        <Legend />
        <Bar
          dataKey="before"
          fill={theme.colors.hexToRgba(theme.colors.strongNegative, 0.4)}
          barSize={barSize}
          shape={<CustomBarShape dataKey="before" />}
        />
        <Bar
          dataKey="after"
          fill={theme.colors.hexToRgba(theme.colors.strongPositive, 0.4)}
          barSize={barSize}
          shape={<CustomBarShape dataKey="after" />}
        />
      </BarChart>
    </ResponsiveContainer>
  );
};

export default DualMetricBarChart;
