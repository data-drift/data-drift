import type { Meta, StoryObj } from "@storybook/react";

import { StepChart } from "./StepChart";

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
  args: {
    metricNames: ["2023-02", "2022-01"],

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
      { daysSinceFirstReport: 344.423287037037, "2022-01": 0 },
      { daysSinceFirstReport: 344.4266087962963, "2022-01": 0 },
      { daysSinceFirstReport: 344.4275, "2022-01": 0 },
      { daysSinceFirstReport: 344.42751157407406, "2022-01": 0 },
      { daysSinceFirstReport: 344.4289814814815, "2022-01": 0 },
      { daysSinceFirstReport: 344.42899305555557, "2022-01": 0 },
      { daysSinceFirstReport: 344.48516203703707, "2022-01": 0 },
      { daysSinceFirstReport: 345.22763888888886, "2022-01": 0 },
      { daysSinceFirstReport: 346.579525462963, "2022-01": 0 },
      { daysSinceFirstReport: 350.2256597222222, "2022-01": 0 },
      { daysSinceFirstReport: 357.4198263888889, "2022-01": 0 },
      { daysSinceFirstReport: 364.2250925925926, "2022-01": 0 },
      { daysSinceFirstReport: 365.47398148148153, "2022-01": 0 },
      { daysSinceFirstReport: 365.47414351851853, "2022-01": 0 },
      { daysSinceFirstReport: 371.22749999999996, "2022-01": 0 },
      { daysSinceFirstReport: 372.23555555555555, "2022-01": 0 },
      { daysSinceFirstReport: 373.6502083333333, "2022-01": 0 },
      { daysSinceFirstReport: 373.7151157407407, "2022-01": 0 },
      { daysSinceFirstReport: 379.44136574074076, "2022-01": 0.03263395504281 },
      { daysSinceFirstReport: 409.48615740740746, "2022-01": 0.45034857959079 },
      { daysSinceFirstReport: 410.230462962963, "2022-01": 0.03263395504281 },
    ],
  },
};
