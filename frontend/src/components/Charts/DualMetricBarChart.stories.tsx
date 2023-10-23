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
      { name: "MRR 2023-01", before: 100, after: 95, percentageChange: -5 },
      { name: "MRR 2023-02", before: 85, after: 95, percentageChange: 11.76 },
      { name: "MRR 2023-03", before: 90, after: 92, percentageChange: 2.22 },
      { name: "MRR 2023-04", before: 88, after: 96, percentageChange: 9.09 },
      { name: "MRR 2023-05", before: 92, after: 94, percentageChange: 2.17 },
      { name: "MRR 2023-06", before: 87, after: 93, percentageChange: 6.9 },
      { name: "MRR 2023-07", before: 89, after: 91, percentageChange: 2.25 },
      { name: "MRR 2023-08", before: 83, after: 90, percentageChange: 8.43 },
      { name: "MRR 2023-09", before: 84, after: 88, percentageChange: 4.76 },
      { name: "MRR 2023-10", before: 86, after: 89, percentageChange: 3.49 },
      { name: "MRR 2023-11", before: 82, after: 85, percentageChange: 3.66 },
      { name: "MRR 2023-12", before: 81, after: 84, percentageChange: 3.7 },
    ],
  },
};

export const QuarterlyMetric: Story = {
  args: {
    data: [
      {
        name: "MRR Q1 2023",
        before: 123395.76,
        after: 123342.12,
        percentageChange: -0.04,
      },
      {
        name: "MRR Q2 2023",
        before: 23295.76,
        after: 57642.12,
        percentageChange: 147.98,
      },
      {
        name: "MRR Q3 2023",
        before: 101395.76,
        after: 123342.12,
        percentageChange: 21.62,
      },
      {
        name: "MRR Q4 2023",
        before: 105395.76,
        after: 105395.76,
        percentageChange: 0,
      },
    ],
  },
};
