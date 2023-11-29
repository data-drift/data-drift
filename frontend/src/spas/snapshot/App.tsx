import { DualTable } from "../../components/Table/DualTable";
import {
  mapSnapshotDiffToDualTableProps,
  sampleSnapshotDiff,
} from "./generatedDiff";

const diff = window.generated_diff || sampleSnapshotDiff;

const tableProps = mapSnapshotDiffToDualTableProps(diff);

const App = () => {
  return (
    <DualTable
      tableProps1={tableProps.tableProps1}
      tableProps2={tableProps.tableProps2}
    />
  );
};

export default App;
