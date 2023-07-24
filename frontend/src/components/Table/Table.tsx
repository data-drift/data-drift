import styled from "@emotion/styled";

// define the CSS

const StyledTd = styled.td`
  width: 300px; // same width for every cell
  height: 50px; // same height for every cell
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
  <table>
    <thead>
      <tr>
        {headers.map((header, i) => (
          <th key={"header" + i}>{header}</th>
        ))}
      </tr>
    </thead>
    <tbody>
      {data.map((row, i) => (
        <tr key={i}>
          {row.map((cell, j) => (
            <StyledTd key={"cell" + i + "-" + j}>{cell}</StyledTd>
          ))}
        </tr>
      ))}
    </tbody>
  </table>
);
