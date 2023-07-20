import styled from "@emotion/styled";
import { useRef } from "react";

const TableContainer = styled.div`
  display: flex;
  justify-content: space-between;
  overflow-x: auto;
`;

const Table = styled.table`
  width: 50%;
  border-collapse: collapse;
  overflow-y: hidden;
  white-space: nowrap;
`;

const TableHeader = styled.th`
  padding: 8px;
  text-align: left;
  border-bottom: 1px solid #ddd;
`;

const LeftCell = styled.td`
  color: red;
  font-weight: bold;
`;

const RightCell = styled.td`
  color: green;
  font-weight: bold;
`;

interface DiffTableProps {
  /**
   * How many lines
   */
  lineCount: number;
  /**
   * How many columns
   */
  headerCount: number;
}

export const DiffTable = ({ lineCount, headerCount }: DiffTableProps) => {
  const leftTableRef = useRef<HTMLTableElement>(null);
  const rightTableRef = useRef<HTMLTableElement>(null);

  const handleScroll = (event: any) => {
    const { scrollTop } = event.target;
    if (leftTableRef.current && rightTableRef.current) {
      leftTableRef.current.scrollTop = scrollTop;
      rightTableRef.current.scrollTop = scrollTop;
    }
    const { scrollLeft } = event.currentTarget;
    if (leftTableRef.current) leftTableRef.current.scrollLeft = scrollLeft;
    if (rightTableRef.current) rightTableRef.current.scrollLeft = scrollLeft;
  };

  return (
    <TableContainer onScroll={handleScroll}>
      <Table ref={leftTableRef}>
        <thead>
          <tr>
            {Array.from({ length: headerCount }).map((_, index) => (
              <TableHeader key={index}>Header {index + 1}</TableHeader>
            ))}
          </tr>
        </thead>
        <tbody>
          {Array.from({ length: lineCount }).map((_, indexValue) => (
            <tr key={indexValue}>
              {Array.from({ length: headerCount }).map((_, index) => (
                <LeftCell>
                  {"ancient"} {indexValue + 1} {index + 1}
                </LeftCell>
              ))}
            </tr>
          ))}
        </tbody>
      </Table>
      <Table ref={rightTableRef} onScroll={handleScroll}>
        <thead>
          <tr>
            {Array.from({ length: headerCount }).map((_, index) => (
              <TableHeader key={index}>Header {index + 1}</TableHeader>
            ))}
          </tr>
        </thead>
        <tbody>
          {Array.from({ length: lineCount }).map((_, indexValue) => (
            <tr key={indexValue}>
              {Array.from({ length: headerCount }).map((_, index) => (
                <RightCell>
                  {"coucou"} {indexValue + 1} {index + 1}
                </RightCell>
              ))}
            </tr>
          ))}
        </tbody>
      </Table>
    </TableContainer>
  );
};
