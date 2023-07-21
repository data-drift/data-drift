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
}

export const Table: React.FC<TableProps> = ({ data }) => (
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
