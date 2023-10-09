import { Position } from "reactflow";

export const nodes = [
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
    position: { x: 50, y: 10 },
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
    position: { x: 450, y: 10 },
    sourcePosition: Position.Right,
    targetPosition: Position.Left,
  },
];
export const edges = [{ id: "e1-2", source: "1", target: "2", animated: true }];

export const mockedDiffTable = {
  tableProps1: {
    diffType: "removed",
    data: Array.from({ length: 130 }).map((_, i) => ({
      isEmphasized: i % 5 === 4,
      data: Array.from({ length: 10 }).map((_, j) => ({
        isEmphasized: i % 5 === 4 && (j + 2 * i) % 6 === 2,
        value: `Old ${i}-${j}`,
      })),
    })),
    headers: Array.from({ length: 10 }).map((_, j) => `Header ${j}`),
  },
  tableProps2: {
    diffType: "added",
    data: Array.from({ length: 130 }).map((_, i) => ({
      isEmphasized: i % 5 === 4,
      data: Array.from({ length: 10 }).map((_, j) => ({
        isEmphasized: i % 5 === 4 && (j + 2 * i) % 6 === 2,
        value: `New ${i}-${j}`,
      })),
    })),
    headers: Array.from({ length: 10 }).map((_, j) => `Header ${j}`),
  },
} as const;
