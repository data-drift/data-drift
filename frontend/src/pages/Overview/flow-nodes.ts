import { Edge, Node, Position } from "reactflow";
import { DDConfigMetric } from "../../services/data-drift";
import { extractFileNameAndPath } from "../../services/string-helpers";
import { LineageEvent } from "../../components/Lineage/MetricNode";

const baseNode = {
  type: "metricNode",
  sourcePosition: Position.Right,
  targetPosition: Position.Left,
};

export const getFileCommits = (
  commitList: {
    commit: { message: string; author: { date?: string } | null };
    sha: string;
  }[],
  filepath: DDConfigMetric["filepath"],
  selectCommit: (commit: string) => void
): LineageEvent[] => {
  filepath = filepath.replace(/\.csv/g, "");
  const metricCommits = commitList.filter((commit) => {
    return commit.commit.message.includes(filepath);
  });
  const metricEvents = metricCommits.reduce((acc, commit) => {
    const isDrift = commit.commit.message.toLowerCase().includes("drift");
    const type = isDrift ? "Drift" : "New Data";
    const eventDate = commit.commit.author?.date
      ? new Date(commit.commit.author.date)
      : null;
    const isPartition = commit.commit.message.includes(filepath + "/");
    if (isPartition) {
      const partitionName = commit.commit.message
        .split("/")
        .pop()
        ?.split(".")[0] as string;
      const subEvent = {
        name: partitionName,
        onClick: () => selectCommit(commit.sha),
      } satisfies NonNullable<LineageEvent["subEvents"]>[0];
      const existingEvent = acc.find((event) => event.type === type);
      if (existingEvent) {
        existingEvent.subEvents?.push(subEvent);
      } else {
        acc.push({
          type,
          eventDate,
          subEvents: [subEvent],
        });
      }
      return acc;
    } else {
      acc.push({
        type,
        eventDate,
        onClick: () => selectCommit(commit.sha),
      });
      return acc;
    }
  }, [] as LineageEvent[]) satisfies LineageEvent[];
  return metricEvents;
};

export const getNodesFromConfig = (
  metric: DDConfigMetric,
  commitList: {
    commit: { message: string; author: { date?: string } | null };
    sha: string;
  }[],
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
  let upstreamNodesOffset = 0;
  const upstreamNodes = metric.upstreamFiles
    ? metric.upstreamFiles.map((upstreamMetric, i) => {
        const upstreamEvents = getFileCommits(
          commitList,
          upstreamMetric,
          selectCommit
        );
        const upstreamNode: Node = {
          ...baseNode,
          position: { x: 50, y: 10 + i * 110 + upstreamNodesOffset },
          id: `upstream-${i}`,
          data: {
            label: extractFileNameAndPath(upstreamMetric).fileName,
            events: upstreamEvents,
          },
        } satisfies Node;
        upstreamNodesOffset += 30 * upstreamEvents.length;

        return upstreamNode;
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
