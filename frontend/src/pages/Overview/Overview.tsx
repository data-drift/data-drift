import Lineage from "../../components/Lineage/Lineage";
import DualTableHeader from "../../components/Table/DualTableHeader";
import { DualTable } from "../../components/Table/DualTable";
import { useCallback, useEffect, useState } from "react";
import { Edge, Node } from "reactflow";

import {
  Container,
  DiffTableContainer,
  LineageContainer,
  StyledCollapsibleContent,
  StyledCollapsibleTitle,
  StyledDate,
  StyledDateButton,
  StyledHeader,
  StyledSelect,
} from "./components";
import { mockedDiffTable } from "./mocked-data";
import { loader, useOverviewLoaderData } from "./loader";
import { getNodesFromConfig } from "./flow-nodes";
import { getCommitList } from "../../services/data-drift";
import { Endpoints } from "@octokit/types";

const Overview = () => {
  const config = useOverviewLoaderData();
  const searchParams = new URLSearchParams(window.location.search);

  const dualTableHeaderState = DualTableHeader.useState();
  const initialSelectedMetric = Number(searchParams.get("metric")) || 0;
  const [selectedMetric, setSelectedMetric] = useState(
    config.config.metrics[initialSelectedMetric]
  );
  const handleSetSelectedMetric = useCallback(
    (newMetricIndex: number) => {
      const searchParams = new URLSearchParams(window.location.search);
      searchParams.set("metric", newMetricIndex.toString());
      const newUrl = `${window.location.pathname}?${searchParams.toString()}`;
      window.history.pushState({ path: newUrl }, "", newUrl);
      setSelectedMetric(config.config.metrics[newMetricIndex]);
    },
    [config.config.metrics]
  );

  const initialSnapshotDate = searchParams.get("snapshotDate")
    ? new Date(searchParams.get("snapshotDate") as string)
    : new Date();
  const [currentDate, setCurrentDate] = useState(initialSnapshotDate);
  const handleSetCurrentDate = useCallback(
    (newDate: Date) => {
      const searchParams = new URLSearchParams(window.location.search);
      searchParams.set(
        "snapshotDate",
        currentDate.toISOString().substring(0, 10)
      );
      const newUrl = `${window.location.pathname}?${searchParams.toString()}`;
      window.history.pushState({ path: newUrl }, "", newUrl);
      setCurrentDate(newDate);
    },
    [currentDate]
  );

  const [commitListData, setCommitListData] = useState({
    data: [] as Endpoints["GET /repos/{owner}/{repo}/commits"]["response"]["data"],
    loading: true,
    nodes: [] as Node[],
    edges: [] as Edge[],
  });

  useEffect(() => {
    const fetchCommit = async () => {
      const result = await getCommitList(
        config.params,
        currentDate.toISOString().substring(0, 10)
      );
      const { nodes, edges } = getNodesFromConfig(selectedMetric, result.data);
      setCommitListData({ data: result.data, loading: false, nodes, edges });
    };
    void fetchCommit();
  }, [currentDate, config.params, selectedMetric]);

  const [isCollapsed, setIsCollapsed] = useState(false);

  const incrementDate = useCallback(() => {
    const newDate = new Date(currentDate.setDate(currentDate.getDate() + 1));
    handleSetCurrentDate(newDate);
  }, [handleSetCurrentDate, currentDate]);

  const decrementDate = useCallback(() => {
    const newDate = new Date(currentDate.setDate(currentDate.getDate() - 1));
    handleSetCurrentDate(newDate);
  }, [handleSetCurrentDate, currentDate]);

  return (
    <Container>
      <StyledHeader>
        <StyledCollapsibleTitle onClick={() => setIsCollapsed(!isCollapsed)}>
          {isCollapsed ? "▶" : "▼"} Lineage
        </StyledCollapsibleTitle>
        <StyledDate>
          <StyledDateButton onClick={decrementDate}>{"<"}</StyledDateButton>
          {currentDate.toLocaleDateString()}
          <StyledDateButton onClick={incrementDate}>{">"}</StyledDateButton>
        </StyledDate>

        <StyledSelect
          value={selectedMetric.filepath}
          onChange={(e) => {
            const selectedMetric = config.config.metrics.findIndex(
              (metric) => metric.filepath === e.target.value
            );
            if (typeof selectedMetric === "number" && !isNaN(selectedMetric)) {
              handleSetSelectedMetric(selectedMetric);
            }
          }}
        >
          {config.config.metrics.map((metric) => (
            <option key={metric.filepath} value={metric.filepath}>
              {metric.filepath}
            </option>
          ))}
        </StyledSelect>
      </StyledHeader>

      <LineageContainer>
        {!isCollapsed && (
          <StyledCollapsibleContent isCollapsed={isCollapsed}>
            <Lineage
              nodes={commitListData.nodes}
              edges={commitListData.edges}
            />
          </StyledCollapsibleContent>
        )}
      </LineageContainer>

      <DiffTableContainer>
        <DualTableHeader
          state={dualTableHeaderState}
          copyAction={console.log}
        />
        <DualTable {...mockedDiffTable} />
      </DiffTableContainer>
    </Container>
  );
};

Overview.loader = loader;

export default Overview;
