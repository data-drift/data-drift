import { useMemo } from "react";
import { DualTable, DualTableProps } from "../../components/Table/DualTable";
import DualTableHeader from "../../components/Table/DualTableHeader";

const filterDualTablePropsData = (
  dualTableProps: DualTableProps,
  startDate: string,
  endDate: string
): DualTableProps => {
  const firstOccurence = dualTableProps.tableProps1.data.findIndex((row) => {
    return row.data.length > 0 && row.data[1].value >= startDate;
  });

  const lastOccurence =
    dualTableProps.tableProps1.data.findIndex((row) => {
      return row.data.length > 0 && row.data[1].value >= endDate;
    }) - 1;

  return {
    tableProps1: {
      ...dualTableProps.tableProps1,
      data: dualTableProps.tableProps1.data.slice(
        firstOccurence,
        lastOccurence
      ),
    },
    tableProps2: {
      ...dualTableProps.tableProps2,
      data: dualTableProps.tableProps2.data.slice(
        firstOccurence,
        lastOccurence
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
