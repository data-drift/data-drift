import styled from "@emotion/styled";

const TableContainer = styled.div`
  display: flex;
  justify-content: space-between;
`;

const Table = styled.table`
  width: 100%;
  border-collapse: collapse;
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
   * What is the ancient data
   */
  ancientData: string[];
  /**
   * What is the new data
   */
  newData: string[];
}

export const DiffTable = ({ ancientData, newData }: DiffTableProps) => {
  return (
    <TableContainer>
      <Table>
        <thead>
          <tr>
            <TableHeader>Header 1</TableHeader>
            <TableHeader>Header 2</TableHeader>
            {/* Add more header columns as needed */}
          </tr>
        </thead>
        <tbody>
          {ancientData.map((ancientValue, index) => (
            <tr key={index}>
              <LeftCell>{ancientValue}</LeftCell>
              {/* Add more left-side cells as needed */}
            </tr>
          ))}
        </tbody>
      </Table>
      <Table>
        <thead>
          <tr>
            <TableHeader>Header 1</TableHeader>
            <TableHeader>Header 2</TableHeader>
            {/* Add more header columns as needed */}
          </tr>
        </thead>
        <tbody>
          {newData.map((newValue, index) => (
            <tr key={index}>
              <RightCell>{newValue}</RightCell>
              {/* Add more right-side cells as needed */}
            </tr>
          ))}
        </tbody>
      </Table>
    </TableContainer>
  );
};
