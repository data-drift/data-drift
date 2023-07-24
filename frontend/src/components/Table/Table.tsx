import styled from "@emotion/styled";

const StyledTable = styled.table`
  background-color: ${(props) => props.theme.colors.background};
  table-layout: fixed;
  border-collapse: collapse;
`;

const StyledTHead = styled.thead`
  color: ${(props) => props.theme.colors.text};
`;

const StyledTr = styled.tr`
  display: flex;
  border-collapse: collapse;
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
  font-size: 16px;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
`;

const StyledTd = styled.td`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  align-self: stretch;
  width: 100%;
  box-sizing: border-box;
  border-collapse: collapse;

  color: ${(props) => props.theme.colors.text2};
  font-size: 10px;
  font-style: normal;
  font-weight: 400;
  line-height: 150%;
  letter-spacing: 0.1px;
  border: 1px solid red;
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
        <StyledTr key={`row-i`}>
          {row.map((cell, j) => (
            <StyledTd key={`cell-${i}-${j}`}>{cell}</StyledTd>
          ))}
        </StyledTr>
      ))}
    </tbody>
  </StyledTable>
);
