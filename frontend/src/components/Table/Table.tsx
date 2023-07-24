import styled from "@emotion/styled";

const StyledTable = styled.table`
  background-color: ${(props) => props.theme.colors.background};
  table-layout: fixed;
  border-collapse: collapse;
  box-sizing: border-box;
`;

const StyledTHead = styled.thead`
  color: ${(props) => props.theme.colors.text};
  background-color: ${(props) => props.theme.colors.background2};
`;

const StyledTr = styled.tr``;

const StyledTh = styled.th`
  // layout
  border: 1px solid ${(props) => props.theme.colors.text};
  width: 100%;
  padding: var(--vertical-padding) var(--horizontal-padding);
  --vertical-padding: ${({ theme }) => theme.spacing(2)};
  --horizontal-padding: ${({ theme }) => theme.spacing(6)};
  white-space: nowrap;

  // text
  font-style: normal;
  font-size: ${(props) => props.theme.text.fontSize.medium};
  font-weight: 500;
  line-height: normal;
`;

const StyledTBody = styled.tbody``;

const StyledTd = styled.td<{
  diffType: TableProps["diffType"];
  isEmphasized: Datum["isEmphasized"];
}>`
  // layout
  width: 100%;
  max-width: 350px;
  padding: var(--vertical-padding) var(--horizontal-padding);
  --vertical-padding: ${({ theme }) => theme.spacing(2)};
  --horizontal-padding: ${({ theme }) => theme.spacing(6)};
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;

  border: 1px solid ${(props) => props.theme.colors.background2};

  // text
  color: ${(props) => props.theme.colors.text2};
  font-style: normal;
  font-size: ${(props) => props.theme.text.fontSize.small};
  font-weight: ${(props) => (props.isEmphasized ? 700 : 400)};
  line-height: 150%;
  letter-spacing: 0.1px;
`;

export interface Datum {
  isEmphasized?: boolean;
  value: string;
}

export interface TableProps {
  // Are those data removed or added
  diffType: "removed" | "added";
  // What are the data to display
  data: Datum[][];
  // What are the headers
  headers: string[];
}

export const Table: React.FC<TableProps> = ({ data, headers, diffType }) => (
  <StyledTable>
    <StyledTHead>
      <StyledTr>
        {headers.map((header, i) => (
          <StyledTh key={`header-${i}`}>{header}</StyledTh>
        ))}
      </StyledTr>
    </StyledTHead>
    <StyledTBody>
      {data.map((row, i) => (
        <StyledTr key={`row-i`}>
          {row.map((cell, j) => (
            <StyledTd
              key={`cell-${i}-${j}`}
              diffType={diffType}
              isEmphasized={cell.isEmphasized}
            >
              {cell.value}
            </StyledTd>
          ))}
        </StyledTr>
      ))}
    </StyledTBody>
  </StyledTable>
);
