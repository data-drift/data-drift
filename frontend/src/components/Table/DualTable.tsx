import React, { useRef, useEffect } from "react";
import styled from "@emotion/styled";

// define the CSS
const TableContainer = styled.div`
  display: inline-block;
  overflow: auto;
  width: 50%;
`;

const StyledTd = styled.td`
  width: 200px; // same width for every cell
  height: 50px; // same height for every cell
`;

interface TableProps {
  data: any[][];
  onScroll: (event: React.UIEvent<HTMLDivElement, UIEvent>) => void;
}

const Table: React.FC<TableProps> = ({ data, onScroll }) => (
  <TableContainer onScroll={onScroll}>
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
  </TableContainer>
);

interface DualTableProps {
  data1: any[][];
  data2: any[][];
}

export const DualTable = ({ data1, data2 }: DualTableProps) => {
  const table1Ref = useRef<HTMLDivElement>(null);
  const table2Ref = useRef<HTMLDivElement>(null);

  const handleScroll =
    (
      src: React.RefObject<HTMLDivElement>,
      dest: React.RefObject<HTMLDivElement>
    ) =>
    (event: React.UIEvent<HTMLDivElement, UIEvent>) => {
      if (dest.current && src.current) {
        dest.current.scrollTop = src.current.scrollTop;
        dest.current.scrollLeft = src.current.scrollLeft;
      }
    };

  useEffect(() => {
    const table1 = table1Ref.current;
    const table2 = table2Ref.current;
    if (table1 && table2) {
      table1.onscroll = handleScroll(table1Ref, table2Ref);
      table2.onscroll = handleScroll(table2Ref, table1Ref);
    }
  }, []);

  return (
    <>
      <Table data={data1} onScroll={handleScroll(table1Ref, table2Ref)} />
      <Table data={data2} onScroll={handleScroll(table2Ref, table1Ref)} />
    </>
  );
};
