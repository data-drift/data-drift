import { Meta, StoryObj } from "@storybook/react";
import Lineage from "./Lineage";
import { Position } from "reactflow";

const containerStyles = {
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
  height: "80vh",
  width: "90vw",
  overflow: "auto",
  border: "1px solid #e0e0e0",
  margin: "1em",
  padding: "1em",
};

const meta = {
  title: "Lineage/Lineage",
  component: Lineage,
  decorators: [
    (Story) => (
      <div style={containerStyles}>
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof Lineage>;

export default meta;

type Story = StoryObj<typeof meta>;

export const DefaultCase: Story = {
  args: {
    nodes: [
      {
        id: "1",
        data: {
          label: "organisation_bop_eop_mrr",
          events: [
            { type: "New Data" },
            {
              type: "Drift",
              subEvents: [{ name: "2023-09" }, { name: "2023-10" }],
            },
          ],
        },
        type: "metricNode",
        position: { x: 50, y: 100 },
        sourcePosition: Position.Right,
        targetPosition: Position.Left,
      },
      {
        id: "2",
        type: "metricNode",
        data: {
          label: "bop_eop_mrr_monthly_by_country",
          events: [{ type: "New Data" }, { type: "Drift" }],
        },
        position: { x: 450, y: 100 },
        sourcePosition: Position.Right,
        targetPosition: Position.Left,
      },
    ],
    edges: [{ id: "e1-2", source: "1", target: "2", animated: true }],
  },
};

export const MultipleParentCase: Story = {
  args: {
    nodes: [
      {
        id: "1",
        data: {
          label: "parent_1",
          events: [
            { type: "New Data" },
            {
              type: "Drift",
              subEvents: [{ name: "2023-09" }, { name: "2023-10" }],
            },
          ],
        },
        type: "metricNode",
        position: { x: 50, y: 100 },
        sourcePosition: Position.Right,
        targetPosition: Position.Left,
      },
      {
        id: "2",
        data: {
          label: "parent_2",
          events: [
            { type: "New Data" },
            {
              type: "Drift",
              subEvents: [{ name: "2023-09" }, { name: "2023-10" }],
            },
          ],
        },
        type: "metricNode",
        position: { x: 50, y: 200 },
        sourcePosition: Position.Right,
        targetPosition: Position.Left,
      },
      {
        id: "3",
        data: {
          label: "parent_3",
          events: [
            { type: "New Data" },
            {
              type: "Drift",
              subEvents: [{ name: "2023-09" }, { name: "2023-10" }],
            },
          ],
        },
        type: "metricNode",
        position: { x: 50, y: 300 },
        sourcePosition: Position.Right,
        targetPosition: Position.Left,
      },
      {
        id: "4",
        type: "metricNode",
        data: {
          label: "child_4",
          events: [{ type: "New Data" }, { type: "Drift" }],
        },
        position: { x: 450, y: 100 },
        sourcePosition: Position.Right,
        targetPosition: Position.Left,
      },
    ],
    edges: [
      { id: "e1-4", source: "1", target: "4", animated: true },
      { id: "e2-4", source: "2", target: "4", animated: true },
      { id: "e3-4", source: "3", target: "4", animated: true },
    ],
  },
};
