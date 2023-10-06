import { ComponentType, useMemo } from "react";
import ReactFlow, { Node, Edge, Position, Handle, NodeProps } from "reactflow";

import "reactflow/dist/style.css";
import MetricNode from "./MetricNode";
import type { LineageEvent } from "./MetricNode";

const CustomMetricNode: ComponentType<
  NodeProps<{ label: string; events: LineageEvent[] }>
> = ({ data, sourcePosition, targetPosition }) => {
  return (
    <>
      {targetPosition && <Handle type="target" position={targetPosition} />}
      <MetricNode metricName={data.label} events={data.events} />

      {sourcePosition && <Handle type="source" position={sourcePosition} />}
    </>
  );
};

const initialNodes = [
  {
    id: "1",
    data: {
      label: "organisation_bop_eop_mrr",
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
    },
    position: { x: 450, y: 100 },
    sourcePosition: Position.Right,
    targetPosition: Position.Left,
  },
] satisfies Node[];

const initialEdges = [
  { id: "e1-2", source: "1", target: "2", animated: true },
] satisfies Edge[];

function Flow() {
  const nodeTypes = useMemo(() => ({ metricNode: CustomMetricNode }), []);

  return (
    <div style={{ width: "1000px", height: "1000px" }}>
      <ReactFlow
        draggable={false}
        nodes={initialNodes}
        edges={initialEdges}
        style={{ width: "600px", height: "300px" }}
        nodesDraggable={false}
        edgesUpdatable={false}
        nodeTypes={nodeTypes}
      ></ReactFlow>
    </div>
  );
}

export default Flow;
