import { Edge, Node, Position } from "reactflow";
import { nodes, edges } from "./mocked-data";
import { DDConfigMetric } from "../../services/data-drift";
import { extractFileNameAndPath } from "../../services/string-helpers";

const baseNode = {
  type: "metricNode",
  sourcePosition: Position.Right,
  targetPosition: Position.Left,
};

export const getNodesFromConfig = (
  metric: DDConfigMetric
): { nodes: Node[]; edges: Edge[] } => {
  const metricNode: Node = {
    ...baseNode,
    id: "2",
    position: { x: 450, y: 10 },
    data: {
      label: extractFileNameAndPath(metric.filepath).fileName,
      events: [],
    },
  } satisfies Node;
  const upstreamNodes = metric.upstreamFiles
    ? metric.upstreamFiles.map((upstreamMetric, i) => {
        console.log("upstreamMetric", upstreamMetric);
        return {
          ...baseNode,
          position: { x: 50, y: 10 },
          id: "1",
          data: {
            label: extractFileNameAndPath(upstreamMetric).fileName,
            events: [],
          },
        } satisfies Node;
      })
    : [];
  return { nodes: [metricNode, ...upstreamNodes], edges };
};
