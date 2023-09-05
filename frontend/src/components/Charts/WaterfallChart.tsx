import { BarChart, Bar, XAxis, YAxis, Tooltip } from "recharts";
import { ComponentProps } from "react";

import { theme } from "../../theme";

type TooltipProps = ComponentProps<typeof Tooltip<[number, number], "drift">>;
type Formatter = TooltipProps["formatter"];

interface DataItem {
  day: string;
  drift: readonly [number, number];
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
  const signedDrift = drift > 0 ? `+${drift.toFixed(2)}` : drift.toFixed(2);
  return `${signedDrift} => ${tickValue[1]}`;
};

export type WaterfallChartProps = { data: readonly DataItem[] };

export const WaterfallChart = ({ data }: WaterfallChartProps) => {
  return (
    <BarChart width={750} height={250} data={[...data]}>
      <XAxis dataKey="day" />
      <YAxis type="number" domain={["auto", "auto"]} tickCount={5} />
      <Tooltip
        contentStyle={{ backgroundColor: theme.colors.background }}
        formatter={formatToolTipValueTick}
      />
      <Bar dataKey="drift" fill={"white"} />
    </BarChart>
  );
};
