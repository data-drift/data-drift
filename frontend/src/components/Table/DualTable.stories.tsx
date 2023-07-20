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

const dataarray = Array.from({ length: 100 }).map((_, i) =>
  Array.from({ length: 10 }).map((_, j) => `Old ${i}-${j}`)
);

export const SimpleCase: Story = {
  args: {
    data1: dataarray,
    data2: dataarray,
  },
};
