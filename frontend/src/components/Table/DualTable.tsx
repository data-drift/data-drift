import React, { useRef } from "react";
import styled from "@emotion/styled";
import { Table, TableProps } from "./Table";

const TableContainer = styled.div`
  display: inline-block;
  overflow: auto;
  width: calc(50% - 4px);
  height: 100%;

  &:first-of-type {
    margin-right: 4px;
  }

  &:last-child {
    margin-left: 4px;
  }
`;

const DualTableContainer = styled.div`
  height: 90vh;
  width: 100%;
`;

export interface DualTableProps {
  /**
   * Properties for the left table
   */
  tableProps1: TableProps;
  /**
   * Properties for the right table
   */
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
    <DualTableContainer>
      <TableContainer ref={table1Ref} onScroll={handleScrollLeft}>
        <Table {...tableProps1} />
      </TableContainer>
      <TableContainer ref={table2Ref} onScroll={handleScrollRight}>
        <Table {...tableProps2} />
      </TableContainer>
    </DualTableContainer>
  );
};
