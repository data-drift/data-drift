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

export const SimpleCase: Story = {
  args: {},
};
