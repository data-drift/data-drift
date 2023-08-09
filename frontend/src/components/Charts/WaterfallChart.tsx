import { BarChart, Bar, XAxis, YAxis, Tooltip } from "recharts";
import { theme } from "../../theme";

const getBarColor = ({ color }: { color: string }) => {
  return color;
};

export const WaterfallChart = () => {
  const data = [
    {
      day: "05-01",
      value: [980, 1000],
      fill: theme.colors.text,
    },
    {
      day: "05-02",
      value: [1000, 1050],
      fill: theme.colors.dataUp,
    },
    {
      day: "05-03",
      value: [1050, 1020],
      fill: theme.colors.dataDown,
    },
    {
      day: "05-04",
      value: [1020, 1010],
      fill: theme.colors.dataDown,
    },
  ];

  return (
    <BarChart
      width={730}
      height={250}
      data={data}
      margin={{ top: 20, right: 20, bottom: 20, left: 20 }}
    >
      <XAxis dataKey="day" />
      <YAxis domain={[980, "dataMax"]} />
      <Tooltip />
      <Bar dataKey="value" />
    </BarChart>
  );
};
