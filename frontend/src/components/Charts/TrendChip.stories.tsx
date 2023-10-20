import type { Meta, StoryObj } from "@storybook/react";
import { theme } from "../../theme";
import TrendChip from "./TrendChip";

const meta = {
  title: "Charts/TrendChip",
  component: TrendChip,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
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
    absoluteValue: 2,
  },
};
