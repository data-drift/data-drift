import styled from "@emotion/styled";
import Lineage from "../../components/Lineage/Lineage";
import { Position } from "reactflow";
import DualTableHeader from "../../components/Table/DualTableHeader";
import { DualTable } from "../../components/Table/DualTable";
import { useCallback, useState } from "react";

const Container = styled.div`
  width: 100%;
  box-sizing: border-box;
`;

const StyledHeader = styled.header`
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: ${(props) =>
    props.theme.colors.background2}; // muted background color
  padding: 20px 40px;
  width: 100%;
  box-sizing: border-box;
`;

const StyledDate = styled.div`
  font-size: 32px;
  font-weight: bold;
  color: ${(props) => props.theme.colors.text};
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 4px;
`;

const StyledSelect = styled.select`
  background-color: ${(props) => props.theme.colors.background};
  border: none;
  padding: 10px 20px;
  font-size: 16px;
  cursor: pointer;
`;

const StyledDateButton = styled.button`
  background-color: transparent;
  color: ${(props) => props.theme.colors.text2};
  border: 0;
  border-radius: 0;
  padding: 5px 10px;
  font-size: 18px;
  cursor: pointer;
  &:hover {
    background-color: #bbb; // or any other color indication for hover
  }
`;

const LineageContainer = styled.div`
  background-color: ${(props) => props.theme.colors.background2};
  text-align: left;
`;

const StyledCollapsibleTitle = styled.button`
  cursor: pointer;
  padding: 10px;
  border: none;
  background-color: ${(props) => props.theme.colors.background2};
`;

const StyledCollapsibleContent = styled.div<{ isCollapsed: boolean }>`
  height: ${(props) => (props.isCollapsed ? "0" : "300px")};

  overflow: hidden;
`;

const DiffTableContainer = styled.div``;

const nodes = [
  {
    id: "1",
    data: {
      label: "organisation_bop_eop_mrr",
      events: [
        { type: "New Data" },
        {
          type: "Drift",
          subEvents: [{ name: "2023-09" }, { name: "2023-10" }],
        },
      ],
    },
    type: "metricNode",
    position: { x: 50, y: 10 },
    sourcePosition: Position.Right,
    targetPosition: Position.Left,
  },
  {
    id: "2",
    type: "metricNode",
    data: {
      label: "bop_eop_mrr_monthly_by_country",
      events: [{ type: "New Data" }, { type: "Drift" }],
    },
    position: { x: 450, y: 10 },
    sourcePosition: Position.Right,
    targetPosition: Position.Left,
  },
];
const edges = [{ id: "e1-2", source: "1", target: "2", animated: true }];

export const Overview = () => {
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
        <DualTable
          {...{
            tableProps1: {
              diffType: "removed",
              data: Array.from({ length: 130 }).map((_, i) => ({
                isEmphasized: i % 5 === 4,
                data: Array.from({ length: 10 }).map((_, j) => ({
                  isEmphasized: i % 5 === 4 && (j + 2 * i) % 6 === 2,
                  value: `Old ${i}-${j}`,
                })),
              })),
              headers: Array.from({ length: 10 }).map((_, j) => `Header ${j}`),
            },
            tableProps2: {
              diffType: "added",
              data: Array.from({ length: 130 }).map((_, i) => ({
                isEmphasized: i % 5 === 4,
                data: Array.from({ length: 10 }).map((_, j) => ({
                  isEmphasized: i % 5 === 4 && (j + 2 * i) % 6 === 2,
                  value: `New ${i}-${j}`,
                })),
              })),
              headers: Array.from({ length: 10 }).map((_, j) => `Header ${j}`),
            },
          }}
        />
      </DiffTableContainer>
    </Container>
  );
};
