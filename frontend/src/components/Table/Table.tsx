import styled from "@emotion/styled";
import { theme } from "../../theme";

const StyledTable = styled.table`
  background-color: ${(props) => props.theme.colors.background};
  table-layout: fixed;
  border-collapse: separate;
  border-spacing: 0;
  box-sizing: border-box;
`;

const StyledTHead = styled.thead`
  color: ${(props) => props.theme.colors.text};
  background-color: ${(props) => props.theme.colors.background2};

  position: sticky;
  top: 0;
`;

const StyledTh = styled.th`
  // layout
  border: 0.5px solid ${(props) => props.theme.colors.text};
  border-top: 1px solid;
  border-bottom: 1px solid;
  width: 100%;
  padding: var(--vertical-padding) var(--horizontal-padding);
  --vertical-padding: ${({ theme }) => theme.spacing(2)};
  --horizontal-padding: ${({ theme }) => theme.spacing(6)};
  white-space: nowrap;

  position: sticky;
  top: 0;

  // text
  font-style: normal;
  font-size: ${(props) => props.theme.text.fontSize.medium};
  font-weight: 500;
  line-height: normal;
`;

const StyledTBody = styled.tbody``;

const StyledTr = styled.tr``;

const EllispsisTd = styled.td`
  width: 100%;
  height: ${({ theme }) => theme.spacing(2)};
  background-color: ${(props) => props.theme.colors.background2};
`;

type DataType = "string" | "number" | "boolean";

const StyledTd = styled.td<{
  diffType: TableProps["diffType"];
  isEmphasized: "cell" | "row" | undefined;
  dataType: DataType;
}>`
  // layout
  width: 100%;
  max-width: 350px;
  padding: var(--vertical-padding) var(--horizontal-padding);
  --vertical-padding: ${({ theme }) => theme.spacing(2)};
  --horizontal-padding: ${({ theme }) => theme.spacing(6)};
  text-align: ${({ dataType }) => (dataType === "number" ? "right" : "center")};
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;

  border: 0.5px solid ${(props) => props.theme.colors.background2};

  // text
  color: ${(props) => props.theme.colors.text2};
  font-style: normal;
  font-size: ${(props) => props.theme.text.fontSize.small};
  font-weight: ${(props) => (props.isEmphasized ? 700 : 400)};
  line-height: 150%;
  letter-spacing: 0.1px;

  background-color: ${(props) =>
    props.isEmphasized === "cell"
      ? props.theme.colors.diffSelect(props.diffType).strong
      : props.isEmphasized === "row"
      ? props.theme.colors.diffSelect(props.diffType).light
      : "inherit"};
`;

export interface Datum {
  isEmphasized?: boolean;
  value: string;
  type?: DataType;
  diffValue?: number;
}

export interface Row {
  data: Datum[];
  isEmphasized?: boolean;
  isEmpty?: boolean;
  isEllipsis?: boolean;
}

export interface TableProps {
  // Are those data removed or added
  diffType: "removed" | "added";
  // What are the data to display
  data: Row[];
  // What are the headers
  headers: string[];
}

const Difference = ({
  diffValue,
  diffType,
}: {
  diffValue: number;
  diffType: "removed" | "added";
}) => {
  const opacity = diffType === "removed" ? 0 : 1;
  const color = diffValue > 0 ? "#0B6E99" : "#FFA344";
  return (
    <span style={{ opacity }}>
      (
      <span style={{ color }}>
        {diffValue > 0 ? (
          <>
            ▲{" "}
            <span style={{ color: theme.colors.text }}>
              {diffValue.toLocaleString()}
            </span>
          </>
        ) : (
          <>
            ▼{" "}
            <span style={{ color: theme.colors.text }}>
              {diffValue.toLocaleString()}
            </span>
          </>
        )}
      </span>
      ){" "}
    </span>
  );
};

export const Table: React.FC<TableProps> = ({ data, headers, diffType }) => (
  <StyledTable>
    <colgroup>
      {Array.from({ length: headers.length }).map((_, i) => (
        <col key={`col-${diffType}-${i}`}></col>
      ))}
    </colgroup>
    <StyledTHead>
      <StyledTr>
        {headers.map((header, i) => (
          <StyledTh key={`header-${diffType}-${i}`}>{header}</StyledTh>
        ))}
      </StyledTr>
    </StyledTHead>
    <StyledTBody>
      {data.map((row, i) => (
        <StyledTr key={`row-${diffType}-${i}`}>
          {row.isEllipsis ? (
            <EllispsisTd
              colSpan={headers.length}
              key={`ellipsis-${diffType}-${i}`}
            ></EllispsisTd>
          ) : (
            row.data.map((cell, j) => (
              <StyledTd
                key={`cell-${diffType}-${i}-${j}`}
                diffType={diffType}
                dataType={cell.type || "string"}
                isEmphasized={
                  cell.isEmphasized
                    ? "cell"
                    : row.isEmphasized
                    ? "row"
                    : undefined
                }
                title={cell.value}
              >
                {cell.diffValue && (
                  <Difference diffValue={cell.diffValue} diffType={diffType} />
                )}
                {cell.value}
              </StyledTd>
            ))
          )}
        </StyledTr>
      ))}
    </StyledTBody>
  </StyledTable>
);
