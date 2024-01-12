import { ComponentType, useMemo } from "react";
import ReactFlow, { Node, Edge, Handle, NodeProps } from "reactflow";

import "reactflow/dist/style.css";
import MetricNode from "./MetricNode";
import type { LineageEvent } from "./MetricNode";
import styled from "@emotion/styled";

const CustomMetricNode: ComponentType<
  NodeProps<{ label: string; events: LineageEvent[]; eventsLoading: boolean }>
> = ({ data, sourcePosition, targetPosition }) => {
  return (
    <div>
      {targetPosition && <Handle type="target" position={targetPosition} />}
      <MetricNode
        metricName={data.label}
        events={data.events}
        eventsLoading={data.eventsLoading}
      />

      {sourcePosition && <Handle type="source" position={sourcePosition} />}
    </div>
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
        draggable={true}
        zoomOnScroll={false}
        zoomOnPinch={true}
        zoomOnDoubleClick={false}
        zoomActivationKeyCode={null}
        panOnDrag={true}
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
