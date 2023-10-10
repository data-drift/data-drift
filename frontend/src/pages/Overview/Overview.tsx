import Lineage from "../../components/Lineage/Lineage";
import DualTableHeader from "../../components/Table/DualTableHeader";
import { DualTable } from "../../components/Table/DualTable";
import { useCallback, useState } from "react";
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

const Overview = () => {
  const config = useOverviewLoaderData();
  const dualTableHeaderState = DualTableHeader.useState();

  const [selectedMetric, setSelectedMetric] = useState(
    config.config.metrics[0]
  );
  const [currentDate, setCurrentDate] = useState(new Date());
  const [isCollapsed, setIsCollapsed] = useState(false);

  const incrementDate = useCallback(() => {
    setCurrentDate(
      (prevDate) => new Date(prevDate.setDate(prevDate.getDate() + 1))
    );
  }, []);

  const decrementDate = useCallback(() => {
    setCurrentDate(
      (prevDate) => new Date(prevDate.setDate(prevDate.getDate() - 1))
    );
  }, []);

  const { nodes, edges } = getNodesFromConfig(selectedMetric);
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
            const selectedMetric = config.config.metrics.find(
              (metric) => metric.filepath === e.target.value
            );
            selectedMetric && setSelectedMetric(selectedMetric);
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
            <Lineage nodes={nodes} edges={edges} />
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
