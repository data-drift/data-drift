import Lineage from "../../components/Lineage/Lineage";
import { useCallback, useMemo, useState } from "react";

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
import {
  fetchCommitPatchQuery,
  fetchCommitsQuery,
  loader,
  localStrategyLoader,
  useOverviewLoaderData,
} from "./loader";
import { getNodesFromConfig } from "./flow-nodes";
import { DiffTable } from "../DisplayCommit/DiffTable";
import Loader from "../../components/Common/Loader";
import StarUs from "../../components/Common/StarUs";
import { useQuery } from "@tanstack/react-query";

const Overview = () => {
  const loaderData = useOverviewLoaderData();
  const searchParams = new URLSearchParams(window.location.search);

  const tableName = searchParams.get("tableName") || "";
  const initialSelectedMetricNumber = Number(searchParams.get("metric")) || 0;

  const initialSelectedMetric = useMemo(() => {
    if (tableName.length > 0) {
      const metric = loaderData.config.metrics.find((metric) =>
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
      return loaderData.config.metrics[initialSelectedMetricNumber];
    }
  }, [tableName, loaderData.config.metrics, initialSelectedMetricNumber]);

  const [selectedMetric, setSelectedMetric] = useState(initialSelectedMetric);
  const handleSetSelectedMetric = useCallback(
    (newMetricIndex: number) => {
      const searchParams = new URLSearchParams(window.location.search);
      searchParams.set("metric", newMetricIndex.toString());
      const newUrl = `${window.location.pathname}?${searchParams.toString()}`;
      window.history.pushState({ path: newUrl }, "", newUrl);
      setSelectedMetric(loaderData.config.metrics[newMetricIndex]);
    },
    [loaderData.config.metrics]
  );

  const initialSnapshotDate = searchParams.get("snapshotDate")
    ? new Date(searchParams.get("snapshotDate") as string)
    : new Date();
  const [currentDate, setCurrentDate] = useState(initialSnapshotDate);

  const initialCommitSha = searchParams.get("commitSha");
  const [selectedCommit, setSelectedCommit] = useState(initialCommitSha);
  const handleSetSelectedCommit = useCallback((newCommitSha: string) => {
    const searchParams = new URLSearchParams(window.location.search);
    searchParams.set("commitSha", newCommitSha);
    const newUrl = `${window.location.pathname}?${searchParams.toString()}`;
    window.history.pushState({ path: newUrl }, "", newUrl);
    setSelectedCommit(newCommitSha);
  }, []);

  const dualTableData = useQuery(
    fetchCommitPatchQuery(loaderData, selectedCommit)
  );

  const commitListData = useQuery(fetchCommitsQuery(loaderData, currentDate));

  const { nodes, edges } = getNodesFromConfig(
    selectedMetric,
    commitListData.data || [],
    handleSetSelectedCommit,
    commitListData.isLoading
  );

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

        {loaderData.config.metrics.length > 0 && (
          <StyledSelect
            value={selectedMetric.filepath}
            onChange={(e) => {
              const selectedMetric = loaderData.config.metrics.findIndex(
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
            {loaderData.config.metrics
              .reduce((unique, metric) => {
                return unique.some((item) => item.filepath === metric.filepath)
                  ? unique
                  : [...unique, metric];
              }, [] as typeof loaderData.config.metrics)
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
            <Lineage nodes={nodes} edges={edges} />
          </StyledCollapsibleContent>
        )}
      </LineageContainer>

      {selectedCommit ? (
        <DiffTableContainer>
          {dualTableData.isLoading ? (
            <Loader />
          ) : (
            dualTableData.data && (
              <DiffTable dualTableProps={dualTableData.data} />
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
