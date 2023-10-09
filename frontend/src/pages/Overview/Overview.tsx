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
import { nodes, edges, mockedDiffTable } from "./mocked-data";

const Overview = () => {
  const dualTableHeaderState = DualTableHeader.useState();

  const availableMetrics = ["Metric1", "Metric2", "Metric3"];
  const [selectedMetric, setSelectedMetric] = useState(availableMetrics[0]);
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
          value={selectedMetric}
          onChange={(e) => setSelectedMetric(e.target.value)}
        >
          {availableMetrics.map((metric) => (
            <option key={metric} value={metric}>
              {metric}
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

export default Overview;
