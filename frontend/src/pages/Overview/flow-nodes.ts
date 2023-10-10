import { Edge, Node, Position } from "reactflow";
import { DDConfigMetric, getCommitList } from "../../services/data-drift";
import { extractFileNameAndPath } from "../../services/string-helpers";
import { LineageEvent } from "../../components/Lineage/MetricNode";

const baseNode = {
  type: "metricNode",
  sourcePosition: Position.Right,
  targetPosition: Position.Left,
};

const getFileCommits = (
  commitList: Awaited<ReturnType<typeof getCommitList>>["data"],
  filepath: DDConfigMetric["filepath"],
  selectCommit: (commit: string) => void
): LineageEvent[] => {
  const metricCommits = commitList.filter((commit) => {
    return commit.commit.message.includes(filepath);
  });
  const metricEvents = metricCommits.map((commit) => {
    const isDrift = commit.commit.message.includes("Drift");

    return {
      type: isDrift ? "Drift" : "New Data",
      onClick: () => selectCommit(commit.sha),
    };
  }) satisfies LineageEvent[];
  return metricEvents;
};

export const getNodesFromConfig = (
  metric: DDConfigMetric,
  commitList: Awaited<ReturnType<typeof getCommitList>>["data"],
  selectCommit: (commit: string) => void
): { nodes: Node[]; edges: Edge[] } => {
  const metricEvents = getFileCommits(
    commitList,
    metric.filepath,
    selectCommit
  );
  const metricNode: Node = {
    ...baseNode,
    id: "metric",
    position: { x: 650, y: 10 },
    data: {
      label: extractFileNameAndPath(metric.filepath).fileName,
      events: metricEvents,
    },
  } satisfies Node;
  const upstreamNodes = metric.upstreamFiles
    ? metric.upstreamFiles.map((upstreamMetric, i) => {
        const upstreamEvents = getFileCommits(
          commitList,
          upstreamMetric,
          selectCommit
        );
        return {
          ...baseNode,
          position: { x: 50, y: 10 + i * 100 },
          id: `upstream-${i}`,
          data: {
            label: extractFileNameAndPath(upstreamMetric).fileName,
            events: upstreamEvents,
          },
        } satisfies Node;
      })
    : [];

  const edges = upstreamNodes.map((upstreamNode, i) => ({
    id: `edge-${i}`,
    source: upstreamNode.id,
    target: metricNode.id,
    animated: true,
  }));
  return { nodes: [metricNode, ...upstreamNodes], edges };
};
