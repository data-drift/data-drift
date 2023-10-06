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

type LineageProps = {
  nodes: Node[];
  edges: Edge[];
};

function Lineage({ nodes, edges }: LineageProps) {
  const nodeTypes = useMemo(() => ({ metricNode: CustomMetricNode }), []);

  return (
    <div style={{ width: "1000px", height: "1000px" }}>
      <ReactFlow
        draggable={false}
        nodes={nodes}
        edges={edges}
        style={{ width: "600px", height: "300px" }}
        nodesDraggable={false}
        edgesUpdatable={false}
        nodeTypes={nodeTypes}
      ></ReactFlow>
    </div>
  );
}

export default Lineage;
