import type { Meta, StoryObj } from "@storybook/react";
import { WaterfallChart } from "./WaterfallChart";
import { theme } from "../../theme";

const meta = {
  title: "Charts/WaterfallChart",
  component: WaterfallChart,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
} satisfies Meta<typeof WaterfallChart>;

export default meta;

type Story = StoryObj<typeof meta>;

export const SimpleCase: Story = {
  args: {
    data: [
      {
        day: "05-01",
        drift: [98000000, 102500000],
        fill: theme.colors.text,
        isInitial: true,
      },
      {
        day: "05-02",
        drift: [102500000, 105000000],
        fill: theme.colors.dataUp,
      },
      {
        day: "05-03",
        drift: [105000000, 102000000],
        fill: theme.colors.dataDown,
      },
      {
        day: "05-04",
        drift: [102000000, 101000000],
        fill: theme.colors.dataDown,
      },
    ],
  },
};
