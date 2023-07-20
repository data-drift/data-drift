import type { Meta, StoryObj } from "@storybook/react";

import { DualTable } from "./DualTable";

const meta = {
  title: "Drift/DualTable",
  component: DualTable,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
} satisfies Meta<typeof DualTable>;

export default meta;
type Story = StoryObj<typeof meta>;

export const SimpleCase: Story = {
  args: {
    tableProps1: {
      data: Array.from({ length: 100 }).map((_, i) =>
        Array.from({ length: 10 }).map((_, j) => `Old ${i}-${j}`)
      ),
    },
    tableProps2: {
      data: Array.from({ length: 100 }).map((_, i) =>
        Array.from({ length: 10 }).map((_, j) => `New ${i}-${j}`)
      ),
    },
  },
};
