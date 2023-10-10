import { ComponentType, useMemo } from "react";
import ReactFlow, { Node, Edge, Handle, NodeProps } from "reactflow";

import "reactflow/dist/style.css";
import MetricNode from "./MetricNode";
import type { LineageEvent } from "./MetricNode";
import styled from "@emotion/styled";

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

const StyledContainer = styled.div`
  width: 100%;
  height: 100%;
  .react-flow__renderer,
  .react-flow__renderer * {
    cursor: default !important;
  }
`;

function Lineage({ nodes, edges }: LineageProps) {
  const nodeTypes = useMemo(() => ({ metricNode: CustomMetricNode }), []);

  return (
    <StyledContainer>
      <ReactFlow
        preventScrolling={false}
        draggable={false}
        zoomOnScroll={false}
        zoomOnPinch={false}
        zoomOnDoubleClick={false}
        zoomActivationKeyCode={null}
        panOnDrag={false}
        nodes={nodes}
        edges={edges}
        style={{ width: "600px", height: "300px" }}
        nodesDraggable={false}
        edgesUpdatable={false}
        nodeTypes={nodeTypes}
      ></ReactFlow>
    </StyledContainer>
  );
}

export default Lineage;
