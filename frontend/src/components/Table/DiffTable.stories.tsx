import type { Meta, StoryObj } from "@storybook/react";

import { DiffTable } from "./DiffTable";

const meta = {
  title: "Drift/DiffTable",
  component: DiffTable,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
} satisfies Meta<typeof DiffTable>;

export default meta;
type Story = StoryObj<typeof meta>;

export const SimpleCase: Story = {
  args: {
    ancientData: ["ancient", "data"],
    newData: ["new", "data"],
  },
};
