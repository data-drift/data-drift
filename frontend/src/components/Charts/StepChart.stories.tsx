import type { Meta, StoryObj } from "@storybook/react";

import { StepChart, YearMonthString } from "./StepChart";

import { args, payload } from "./payload";

const meta = {
  title: "Charts/StepChart",
  component: StepChart,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
} satisfies Meta<typeof StepChart>;

export default meta;

type Story = StoryObj<typeof meta>;

export const SingleMetric: Story = {
  args: {
    metricNames: ["2023-02"],
    data: [
      { daysSinceFirstReport: 0.4480324074074074, "2023-02": 0 },
      {
        daysSinceFirstReport: 15.706261574074075,
        "2023-02": -0.17640146440354,
      },
      {
        daysSinceFirstReport: 16.486157407407408,
        "2023-02": -0.11398248469152,
      },
      {
        daysSinceFirstReport: 17.230462962962964,
        "2023-02": -0.17640146440354,
      },
      {
        daysSinceFirstReport: 20.222453703703703,
        "2023-02": -0.11398248469152,
      },
      {
        daysSinceFirstReport: 49.37809027777778,
        "2023-02": 0.31480876724324,
      },
      {
        daysSinceFirstReport: 57.243946759259266,
        "2023-02": 0.44869295561106,
      },
    ],
  },
};

export const MultipleMetrics: Story = {
  args: args,
};
