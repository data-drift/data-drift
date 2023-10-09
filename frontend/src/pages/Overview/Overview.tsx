import Lineage from "../../components/Lineage/Lineage";
import { Position } from "reactflow";
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
