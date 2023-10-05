import ReactFlow, { Node, Edge, Position } from "reactflow";

import "reactflow/dist/style.css";

const initialNodes = [
  {
    id: "1",
    type: "input",
    data: { label: "FCT Order" },
    position: { x: 50, y: 100 },
    sourcePosition: Position.Right,
  },
  {
    id: "2",
    type: "output",
    data: { label: "Statistic" },
    position: { x: 250, y: 100 },
    targetPosition: Position.Left,
  },
] satisfies Node[];

const initialEdges = [
  { id: "e1-2", source: "1", target: "2", animated: true, type: "straight" },
] satisfies Edge[];

function Flow() {
  return (
    <div style={{ width: "1000px", height: "1000px" }}>
      <ReactFlow
        draggable={false}
        nodes={initialNodes}
        edges={initialEdges}
        style={{ width: "600px", height: "300px" }}
        nodesDraggable={false}
        edgesUpdatable={false}
      />
    </div>
  );
}

export default Flow;
