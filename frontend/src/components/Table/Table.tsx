import styled from "@emotion/styled";

const StyledTable = styled.table`
  background-color: ${(props) => props.theme.colors.background};
`;

const StyledTHead = styled.thead`
  color: ${(props) => props.theme.colors.text};
`;

const StyledTr = styled.tr`
  display: flex;
  align-items: flex-start;
  align-self: stretch;
`;

const StyledTh = styled.th`
  display: flex;
  padding: 8px 24px;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 8px;
  align-self: stretch;
  border-right: 1px solid ${(props) => props.theme.colors.text};
`;

const StyledTd = styled.td`
  width: 300px; // same width for every cell
  height: 50px; // same height for every cell
  color: ${(props) => props.theme.colors.text2};
`;

export interface TableProps {
  // Are those data removed or added
  diffType: "removed" | "added";
  // What are the data to display
  data: any[][];
  // What are the headers
  headers: string[];
}

export const Table: React.FC<TableProps> = ({ data, headers }) => (
  <StyledTable>
    <StyledTHead>
      <StyledTr>
        {headers.map((header, i) => (
          <StyledTh key={`header-${i}`}>{header}</StyledTh>
        ))}
      </StyledTr>
    </StyledTHead>
    <tbody>
      {data.map((row, i) => (
        <tr key={`row-i`}>
          {row.map((cell, j) => (
            <StyledTd key={`cell-${i}-${j}`}>{cell}</StyledTd>
          ))}
        </tr>
      ))}
    </tbody>
  </StyledTable>
);
