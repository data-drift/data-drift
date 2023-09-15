import { DualTableProps } from "../components/Table/DualTable";

export function toTsv(props: DualTableProps): string {
  const { tableProps1, tableProps2 } = props;
  const leftHeaders = tableProps1.headers;
  const rightHeaders = tableProps2.headers;
  const headers = [...leftHeaders, ...rightHeaders];
  const rows = tableProps1.data.map((leftRow, i) => {
    const rightRow = tableProps2.data[i];
    const cells = [
      ...leftRow.data.map((cell) => cell.value),
      ...rightRow.data.map((cell) => cell.value),
    ];
    return cells.join("\t");
  });
  const tsv = [headers.join("\t"), ...rows].join("\n");
  return tsv;
}
