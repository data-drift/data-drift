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

export const EmphasizedTable: Story = {
  args: {
    headers: ["header1", "header2", "header3"],
    data: [
      { data: [{ value: "data" }, { value: "data" }, { value: "data" }] },
      {
        data: [
          { value: "data" },
          { value: "new data", isEmphasized: true },
          { value: "data" },
        ],
        isEmphasized: true,
      },
      { data: [{ value: "data" }, { value: "data" }, { value: "data" }] },
    ],
    diffType: "added",
  },
};

export const LongValueTable: Story = {
  args: {
    headers: ["header1", "header2"],
    data: [
      {
        data: [
          {
            value:
              "YYYY-MM-DD bfb63190-cb98-4a28-b165-f1345b01733c 4a28-b165-f1345b01733c",
          },
          { value: "removed data" },
        ],
      },
      {
        data: [
          { value: "removed data" },
          { value: "YYYY-MM-DD bfb63190-cb98-4a28-b165-f1345b01733c" },
        ],
      },
    ],
    diffType: "removed",
  },
};

export const EmptyLineTable: Story = {
  args: {
    headers: ["header1", "header2", "header3"],
    data: [
      { data: [{ value: "data" }, { value: "data" }, { value: "data" }] },
      {
        data: [{ value: "_" }, { value: "_" }, { value: "_" }],
      },
      { data: [{ value: "data" }, { value: "data" }, { value: "data" }] },
    ],
    diffType: "added",
  },
};

export const EllipsisLineTable: Story = {
  args: {
    headers: ["header1", "header2", "header3"],
    data: [
      { data: [{ value: "data" }, { value: "data" }, { value: "data" }] },
      {
        isEllipsis: true,
        data: [],
      },
      { data: [{ value: "data" }, { value: "data" }, { value: "data" }] },
    ],
    diffType: "added",
  },
};
