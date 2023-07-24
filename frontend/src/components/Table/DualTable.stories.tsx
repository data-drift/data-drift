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
      diffType: "removed",
      data: Array.from({ length: 100 }).map((_, i) => ({
        isEmphasized: i % 5 === 4,
        data: Array.from({ length: 10 }).map((_, j) => ({
          isEmphasized: i % 5 === 4 && j % 6 === 2,
          value: `Old ${i}-${j}`,
        })),
      })),
      headers: Array.from({ length: 10 }).map((_, j) => `Header ${j}`),
    },
    tableProps2: {
      diffType: "added",
      data: Array.from({ length: 100 }).map((_, i) => ({
        isEmphasized: i % 5 === 4,
        data: Array.from({ length: 10 }).map((_, j) => ({
          isEmphasized: i % 5 === 4 && j % 6 === 2,
          value: `New ${i}-${j}`,
        })),
      })),
      headers: Array.from({ length: 10 }).map((_, j) => `Header ${j}`),
    },
  },
};
