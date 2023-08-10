import { BarChart, Bar, XAxis, YAxis, Tooltip } from "recharts";
import { ComponentProps } from "react";

// import type { Formatter } from "recharts";
import { theme } from "../../theme";

type TooltipProps = ComponentProps<typeof Tooltip<[number, number], "drift">>;
type Formatter = TooltipProps["formatter"];

interface DataItem {
  day: string;
  drift: [number, number];
  fill: string;
  isInitial?: boolean;
}

const formatToolTipValueTick: Formatter = (
  tickValue,
  _name,
  item: { payload?: DataItem }
) => {
  if (item.payload?.isInitial) return `${tickValue[1]}`;
  const drift = tickValue[1] - tickValue[0];
  const signedDrift = drift > 0 ? `+${drift}` : drift;
  return `${tickValue[1]} (${signedDrift})`;
};

export const WaterfallChart = ({ data }: { data: DataItem[] }) => {
  return (
    <BarChart
      width={730}
      height={250}
      data={data}
      margin={{ top: 20, right: 20, bottom: 20, left: 20 }}
    >
      <XAxis dataKey="day" />
      <YAxis domain={[980, "dataMax"]} />
      <Tooltip
        contentStyle={{ backgroundColor: theme.colors.background }}
        formatter={formatToolTipValueTick}
      />
      <Bar dataKey="drift" fill={"white"} />
    </BarChart>
  );
};