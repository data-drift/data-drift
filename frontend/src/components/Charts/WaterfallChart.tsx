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

export const WaterfallChart = () => {
  const data = [
    {
      day: "05-01",
      drift: [980, 1000],
      fill: theme.colors.text,
      isInitial: true,
    },
    {
      day: "05-02",
      drift: [1000, 1050],
      fill: theme.colors.dataUp,
    },
    {
      day: "05-03",
      drift: [1050, 1020],
      fill: theme.colors.dataDown,
    },
    {
      day: "05-04",
      drift: [1020, 1010],
      fill: theme.colors.dataDown,
    },
  ] satisfies DataItem[];

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
