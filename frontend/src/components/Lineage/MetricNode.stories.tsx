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
    events: [
      { type: "New Data", eventDate: null },
      {
        type: "Drift",
        eventDate: new Date("2021-09-01 12:13:14"),
        subEvents: [{ name: "2023-09" }, { name: "2023-10" }],
      },
    ],
  },
};

export const NoEventCase: Story = {
  args: {
    metricName: "FCT Order",
    events: [],
  },
};
