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
    data: [
      { x: 0.4480324074074074, y: 0 },
      { x: 15.706261574074075, y: -0.17640146440354 },
      { x: 16.486157407407408, y: -0.11398248469152 },
      { x: 17.230462962962964, y: -0.17640146440354 },
      { x: 20.222453703703703, y: -0.11398248469152 },
      { x: 49.37809027777778, y: 0.31480876724324 },
      { x: 57.243946759259266, y: 0.44869295561106 },
    ],
  },
};
