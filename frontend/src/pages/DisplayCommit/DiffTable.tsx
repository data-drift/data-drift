import { DualTable, DualTableProps } from "../../components/Table/DualTable";
import DualTableHeader from "../../components/Table/DualTableHeader";

export const DiffTable = ({ data }: { data: DualTableProps }) => {
  const dualTableHeaderState = DualTableHeader.useState();

  return (
    <>
      <DualTableHeader state={dualTableHeaderState} />
      <DualTable {...data} />
    </>
  );
};
