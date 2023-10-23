import type { Meta, StoryObj } from "@storybook/react";
import TrendChip from "./TrendChip";

const meta = {
  title: "Charts/TrendChip",
  component: TrendChip,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
  decorators: [
    (Story) => (
      <div style={{ width: "128px" }}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof TrendChip>;

export default meta;

type Story = StoryObj<typeof meta>;

export const UpCase: Story = {
  args: {
    trend: "up",
    absoluteValue: 2,
  },
};

export const DownCase: Story = {
  args: {
    trend: "down",
    absoluteValue: 2,
  },
};

export const NeutralCase: Story = {
  args: {
    trend: "neutral",
    absoluteValue: 0,
  },
};

export const LowNumberCase: Story = {
  args: {
    trend: "up",
    absoluteValue: 0.00002,
  },
};

export const LargeNumberCase: Story = {
  args: {
    trend: "down",
    absoluteValue: 11231231,
  },
};
