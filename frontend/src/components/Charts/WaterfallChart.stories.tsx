import type { Meta, StoryObj } from "@storybook/react";
import { WaterfallChart } from "./WaterfallChart";

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

export const SimpleCase: Story = {};
