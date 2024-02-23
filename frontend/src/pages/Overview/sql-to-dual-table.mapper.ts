import { DualTableProps } from "../../components/Table/DualTable";

export const sqlToDualTableMapper = <T extends string, V extends string>(
  oldRows: {
    values: Record<T, string>[];
    columns: T[];
  },
  newRows: {
    values: Record<V, string>[];
    columns: V[];
  }
): DualTableProps => {
  return {
    tableProps1: {
      headers: oldRows.columns,
      diffType: "removed",
      data: oldRows.values.map((row) => {
        return {
          data: oldRows.columns.map((column) => {
            return {
              value: row[column],
            };
          }),
        };
      }),
    },
    tableProps2: {
      headers: newRows.columns,
      diffType: "added",
      data: newRows.values.map((row) => {
        return {
          data: newRows.columns.map((column) => {
            return {
              value: row[column],
            };
          }),
        };
      }),
    },
  };
};
