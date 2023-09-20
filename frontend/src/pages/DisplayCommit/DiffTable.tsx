import { useCallback, useMemo } from "react";
import { DualTable, DualTableProps } from "../../components/Table/DualTable";
import DualTableHeader from "../../components/Table/DualTableHeader";
import { toTsv } from "../../services/tsv";

const filterDualTablePropsData = (
  dualTableProps: DualTableProps,
  startDate: string,
  endDate: string
): DualTableProps => {
  const firstOccurence = dualTableProps.tableProps1.data.findIndex(
    (row, index) => {
      const row2 = dualTableProps.tableProps2.data[index];

      return (
        (row.data.length > 0 &&
          row.data[1].value != "_" &&
          row.data[1].value >= startDate) ||
        (row2.data.length > 0 &&
          row2.data[1].value != "_" &&
          row2.data[1].value >= startDate)
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
    lastOccurence <= 0 ? undefined : lastOccurence - 1;

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

  const copyInClipboard = useCallback(() => {
    const tsvString = toTsv(dualTableProps);
    // Copy to clipboard
    navigator.clipboard
      .writeText(tsvString)
      .then(function () {
        alert("Table copied to clipboard");
      })
      .catch(function (error) {
        alert("Could not copy table to clipboard");
        console.error(error);
      });
  }, [dualTableProps]);

  return (
    <>
      <DualTableHeader
        state={dualTableHeaderState}
        copyAction={copyInClipboard}
      />
      <DualTable {...filteredDualTableProps} />
    </>
  );
};
