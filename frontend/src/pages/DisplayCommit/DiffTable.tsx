import { useMemo } from "react";
import { DualTable, DualTableProps } from "../../components/Table/DualTable";
import DualTableHeader from "../../components/Table/DualTableHeader";

const filterDualTablePropsData = (
  dualTableProps: DualTableProps,
  startDate: string,
  endDate: string
): DualTableProps => {
  console.log(dualTableProps);
  console.log(startDate, endDate);
  const firstOccurence = dualTableProps.tableProps1.data.findIndex(
    (row, index) => {
      const row2 = dualTableProps.tableProps2.data[index];

      return (
        (row.data.length > 0 && row.data[1].value >= startDate) ||
        (row2.data.length > 0 && row2.data[1].value >= startDate)
      );
    }
  );

  const lastOccurence = dualTableProps.tableProps1.data.findIndex(
    (row, index) => {
      const row2 = dualTableProps.tableProps2.data[index];
      return (
        (row.data.length > 0 &&
          row.data[1].value != "_" &&
          row.data[1].value >= endDate) ||
        (row2.data.length > 0 &&
          row2.data[1].value != "_" &&
          row2.data[1].value >= endDate)
      );
    }
  );

  const lastOccurenceWithoutLastLine =
    lastOccurence === -1 ? -1 : lastOccurence - 1;

  console.log("firstOccurence", firstOccurence, "lastOccurence", lastOccurence);

  return {
    tableProps1: {
      ...dualTableProps.tableProps1,
      data: dualTableProps.tableProps1.data.slice(
        firstOccurence,
        lastOccurenceWithoutLastLine
      ),
    },
    tableProps2: {
      ...dualTableProps.tableProps2,
      data: dualTableProps.tableProps2.data.slice(
        firstOccurence,
        lastOccurenceWithoutLastLine
      ),
    },
  };
};

export const DiffTable = ({
  dualTableProps,
}: {
  dualTableProps: DualTableProps;
}) => {
  const dualTableHeaderState = DualTableHeader.useState();

  const filteredDualTableProps = useMemo(
    () =>
      filterDualTablePropsData(
        dualTableProps,
        dualTableHeaderState.startDate,
        dualTableHeaderState.endDate
      ),
    [
      dualTableProps,
      dualTableHeaderState.startDate,
      dualTableHeaderState.endDate,
    ]
  );

  return (
    <>
      <DualTableHeader state={dualTableHeaderState} />
      <DualTable {...filteredDualTableProps} />
    </>
  );
};
