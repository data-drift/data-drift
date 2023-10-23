import type { Meta, StoryObj } from "@storybook/react";

import DualMetricBarChart from "./DualMetricBarChart";

const meta = {
  title: "Charts/DualMetricBarChart",
  component: DualMetricBarChart,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
  decorators: [
    (Story) => (
      <div style={{ width: "1024px" }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof DualMetricBarChart>;

export default meta;

type Story = StoryObj<typeof meta>;

export const MonthlyMetric: Story = {
  args: {
    data: [
      { name: "MRR 2023-01", before: 100, after: 95 },
      { name: "MRR 2023-02", before: 85, after: 95 },
      { name: "MRR 2023-03", before: 90, after: 92 },
      { name: "MRR 2023-04", before: 88, after: 96 },
      { name: "MRR 2023-05", before: 92, after: 94 },
      { name: "MRR 2023-06", before: 87, after: 93 },
      { name: "MRR 2023-07", before: 89, after: 91 },
      { name: "MRR 2023-08", before: 83, after: 90 },
      { name: "MRR 2023-09", before: 84, after: 88 },
      { name: "MRR 2023-10", before: 86, after: 89 },
      { name: "MRR 2023-11", before: 82, after: 85 },
      { name: "MRR 2023-12", before: 81, after: 84 },
    ],
  },
};

export const QuarterlyMetric: Story = {
  args: {
    data: [
      { name: "MRR Q1 2023", before: 123395.76, after: 123342.12 },
      { name: "MRR Q2 2023", before: 23295.76, after: 57642.12 },
      { name: "MRR Q3 2023", before: 101395.76, after: 123342.12 },
      { name: "MRR Q4 2023", before: 105395.76, after: 105395.76 },
    ],
  },
};
