import { Meta, StoryObj } from "@storybook/react";
import MetricNode from "./MetricNode";

const meta = {
  title: "Lineage/MetricNode",
  component: MetricNode,
} satisfies Meta<typeof MetricNode>;

export default meta;

type Story = StoryObj<typeof meta>;

export const DefaultCase: Story = {
  args: {
    metricName: "FCT Order",
    items: ["Drift: 10", "Data Added: 100"],
  },
};
