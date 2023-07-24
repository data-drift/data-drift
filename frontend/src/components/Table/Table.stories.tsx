import type { Meta, StoryObj } from "@storybook/react";
import { Table } from "./Table";

const meta = {
  title: "Drift/Table",
  component: Table,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
} satisfies Meta<typeof Table>;

export default meta;
type Story = StoryObj<typeof meta>;

export const AddedTable: Story = {
  args: {
    headers: ["header1", "header2"],
    data: [
      ["new data", "new data"],
      ["new data", "new data"],
    ],
    diffType: "added",
  },
};

export const RemovedTable: Story = {
  args: {
    headers: ["header1", "header2"],
    data: [
      ["removed data", "removed data"],
      ["removed data", "removed data"],
    ],
    diffType: "removed",
  },
};
