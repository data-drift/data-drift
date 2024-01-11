import Lineage from "../../components/Lineage/Lineage";
import { DualTableProps } from "../../components/Table/DualTable";
import { useCallback, useEffect, useMemo, useState } from "react";
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
import { loader, localStrategyLoader, useOverviewLoaderData } from "./loader";
import { getNodesFromConfig } from "./flow-nodes";
import {
  getCommitList,
  getCommitListLocalStrategy,
  getMeasurement,
  getPatchAndHeader,
} from "../../services/data-drift";
import { DiffTable } from "../DisplayCommit/DiffTable";
import { parsePatch } from "../../services/patch.mapper";
import Loader from "../../components/Common/Loader";
import StarUs from "../../components/Common/StarUs";

const Overview = () => {
  const config = useOverviewLoaderData();
  const searchParams = new URLSearchParams(window.location.search);

  const tableName = searchParams.get("tableName") || "";
  const initialSelectedMetricNumber = Number(searchParams.get("metric")) || 0;

  const initialSelectedMetric = useMemo(() => {
    if (tableName.length > 0) {
      const metric = config.config.metrics.find((metric) =>
        tableName.includes(metric.filepath.replace(".csv", ""))
      );
      return (
        metric || {
          filepath: tableName,
          dateColumnName: "",
          KPIColumnName: "",
          metricName: tableName,
          timeGrains: [],
          dimensions: [],
        }
      );
    } else {
      return config.config.metrics[initialSelectedMetricNumber];
    }
  }, [tableName, config.config.metrics, initialSelectedMetricNumber]);

  const [selectedMetric, setSelectedMetric] = useState(initialSelectedMetric);
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

  const [commitListData, setCommitListData] = useState({
    data: [] as {
      commit: { message: string; author: { date?: string } | null };
      sha: string;
    }[],
    loading: true,
    nodes: [] as Node[],
    edges: [] as Edge[],
  });

  const initialCommitSha = searchParams.get("commitSha");
  const [selectedCommit, setSelectedCommit] = useState(initialCommitSha);
  const handleSetSelectedCommit = useCallback((newCommitSha: string) => {
    const searchParams = new URLSearchParams(window.location.search);
    searchParams.set("commitSha", newCommitSha);
    const newUrl = `${window.location.pathname}?${searchParams.toString()}`;
    window.history.pushState({ path: newUrl }, "", newUrl);
    setSelectedCommit(newCommitSha);
  }, []);

  const [dualTableData, setDualTableData] = useState({
    dualTableProps: undefined as undefined | DualTableProps,
    loading: false,
  });
  useEffect(() => {
    const controller = new AbortController();
    const fetchPatchData = async () => {
      if (!selectedCommit) return;
      switch (config.strategy) {
        case "local": {
          setDualTableData({ dualTableProps: undefined, loading: true });
          const measurementResults = await getMeasurement(
            "default",
            config.params.tableName,
            selectedCommit
          );
          const { oldData, newData } = parsePatch(
            measurementResults.data.Patch,
            measurementResults.data.Headers
          );
          const dualTableProps = {
            tableProps1: oldData,
            tableProps2: newData,
          };
          setDualTableData({ dualTableProps, loading: false });
          break;
        }
        case "github": {
          setDualTableData({ dualTableProps: undefined, loading: true });
          const patchAndHeader = await getPatchAndHeader(
            {
              installationId: config.params.installationId,
              owner: config.params.owner,
              repo: config.params.repo,
              commitSHA: selectedCommit,
            },
            controller
          );
          const { oldData, newData } = parsePatch(
            patchAndHeader.patch,
            patchAndHeader.headers
          );
          const dualTableProps = {
            tableProps1: oldData,
            tableProps2: newData,
          };
          setDualTableData({ dualTableProps, loading: false });
        }
      }
    };
    void fetchPatchData();
    return () => {
      controller.abort();
    };
  }, [selectedCommit, config.params, config.strategy]);

  useEffect(() => {
    const controller = new AbortController();
    const fetchCommit = async () => {
      switch (config.strategy) {
        case "local": {
          const result = await getCommitListLocalStrategy(
            config.params.tableName,
            currentDate.toISOString().substring(0, 10)
          );

          const mappedCommits = result.data.Measurements.map((commit) => ({
            commit: {
              message: commit.Message,
              author: {
                date: commit.Date,
              },
            },
            sha: commit.Sha,
          })) satisfies {
            commit: { message: string; author: { date?: string } | null };
            sha: string;
          }[];

          const { nodes, edges } = getNodesFromConfig(
            selectedMetric,
            mappedCommits,
            handleSetSelectedCommit
          );
          setCommitListData({
            data: mappedCommits,
            loading: false,
            nodes,
            edges,
          });
          break;
        }
        case "github": {
          const result = await getCommitList(
            config.params,
            currentDate.toISOString().substring(0, 10),
            controller
          );
          const { nodes, edges } = getNodesFromConfig(
            selectedMetric,
            result.data,
            handleSetSelectedCommit
          );
          setCommitListData({
            data: result.data,
            loading: false,
            nodes,
            edges,
          });
        }
      }
    };
    void fetchCommit();
    return () => {
      controller.abort();
    };
  }, [
    currentDate,
    config.params,
    config.strategy,
    selectedMetric,
    handleSetSelectedCommit,
  ]);

  const [isCollapsed, setIsCollapsed] = useState(false);

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
      handleSetSelectedCommit("");
      setCommitListData((prev) => ({ ...prev, loading: true }));
      setDualTableData({ dualTableProps: undefined, loading: false });
    },
    [currentDate, handleSetSelectedCommit]
  );

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

        {config.config.metrics.length > 0 && (
          <StyledSelect
            value={selectedMetric.filepath}
            onChange={(e) => {
              const selectedMetric = config.config.metrics.findIndex(
                (metric) => metric.filepath === e.target.value
              );
              if (
                typeof selectedMetric === "number" &&
                !isNaN(selectedMetric)
              ) {
                handleSetSelectedMetric(selectedMetric);
              }
            }}
          >
            {config.config.metrics
              .reduce((unique, metric) => {
                return unique.some((item) => item.filepath === metric.filepath)
                  ? unique
                  : [...unique, metric];
              }, [] as typeof config.config.metrics)
              .map((metric) => (
                <option key={metric.filepath} value={metric.filepath}>
                  {metric.filepath}
                </option>
              ))}
          </StyledSelect>
        )}
        <StarUs />
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

      {selectedCommit ? (
        <DiffTableContainer>
          {dualTableData.loading ? (
            <Loader />
          ) : (
            dualTableData.dualTableProps && (
              <DiffTable dualTableProps={dualTableData.dualTableProps} />
            )
          )}
        </DiffTableContainer>
      ) : (
        "No drift selected"
      )}
    </Container>
  );
};

Overview.loader = loader;
Overview.localStrategyLoader = localStrategyLoader;

export default Overview;
