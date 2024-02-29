import { DualTableProps } from "../../components/Table/DualTable";
import { TableProps } from "../../components/Table/Table";
import { emptyRow } from "../../services/patch.mapper";

export const sqlToDualTableMapper = <T extends string, V extends string>(
  uniqueKeys: string[],
  _oldRows: {
    values: (Record<T, string> & { unique_key: string })[];
    columns: T[];
  },
  _newRows: {
    values: (Record<V, string> & { unique_key: string })[];
    columns: V[];
  },
  initialDualTable: DualTableProps
): DualTableProps => {
  const oldData: TableProps["data"] = [];
  const newData: TableProps["data"] = [];
  uniqueKeys.forEach((uniqueKey: string) => {
    const tablePropsOldRow: TableProps["data"][0] =
      initialDualTable.tableProps1.data.find(
        (row) => row.data[0]?.value === uniqueKey
      ) || emptyRow(initialDualTable.tableProps1.headers.length);

    const tablePropsNewRow: TableProps["data"][0] =
      initialDualTable.tableProps2.data.find(
        (row) => row.data[0]?.value === uniqueKey
      ) || emptyRow(initialDualTable.tableProps2.headers.length);

    oldData.push(tablePropsOldRow);
    newData.push(tablePropsNewRow);
  });

  return {
    tableProps1: {
      headers: initialDualTable.tableProps1.headers,
      diffType: "removed",
      data: oldData,
    },
    tableProps2: {
      headers: initialDualTable.tableProps2.headers,
      diffType: "added",
      data: newData,
    },
  };
};
