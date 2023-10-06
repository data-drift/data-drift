import { Meta, StoryObj } from "@storybook/react";
import Lineage from "./Lineage";
import { Position } from "reactflow";

const meta = {
  title: "Lineage/Lineage",
  component: Lineage,
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
