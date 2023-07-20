import React, { useRef } from "react";
import styled from "@emotion/styled";

// define the CSS
const TableContainer = styled.div`
  display: inline-block;
  overflow: auto;
  width: 50%;
`;

const StyledTd = styled.td`
  width: 300px; // same width for every cell
  height: 50px; // same height for every cell
`;

interface TableProps {
  data: any[][];
}

const Table: React.FC<TableProps> = ({ data }) => (
  <table>
    <thead>
      <tr>
        {data[0].map((_, i) => (
          <th key={i}>Column {i}</th>
        ))}
      </tr>
    </thead>
    <tbody>
      {data.map((row, i) => (
        <tr key={i}>
          {row.map((cell, j) => (
            <StyledTd key={j}>{cell}</StyledTd>
          ))}
        </tr>
      ))}
    </tbody>
  </table>
);

interface DualTableProps {
  tableProps1: TableProps;
  tableProps2: TableProps;
}

export const DualTable = ({ tableProps1, tableProps2 }: DualTableProps) => {
  const table1Ref = useRef<HTMLDivElement>(null);
  const table2Ref = useRef<HTMLDivElement>(null);

  const handleScrollLeft = (
    _scrollEvent: React.UIEvent<HTMLDivElement, UIEvent>
  ) => {
    if (table2Ref.current && table1Ref.current) {
      table2Ref.current.scrollTop = table1Ref.current.scrollTop;
      table2Ref.current.scrollLeft = table1Ref.current.scrollLeft;
    }
  };

  const handleScrollRight = (
    _scrollEvent: React.UIEvent<HTMLDivElement, UIEvent>
  ) => {
    if (table1Ref.current && table2Ref.current) {
      table1Ref.current.scrollTop = table2Ref.current.scrollTop;
      table1Ref.current.scrollLeft = table2Ref.current.scrollLeft;
    }
  };

  return (
    <>
      <TableContainer ref={table1Ref} onScroll={handleScrollLeft}>
        <Table {...tableProps1} />
      </TableContainer>
      <TableContainer ref={table2Ref} onScroll={handleScrollRight}>
        <Table {...tableProps2} />
      </TableContainer>
    </>
  );
};
