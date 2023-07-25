import { useEffect, useState } from "react";
import "./App.css";
import { getCommitFiles } from "./services/github";
import { DualTable, DualTableProps } from "./components/Table/DualTable";
import { parsePatch } from "./services/patch.mapper";

const [owner, repo, commitSHA] = [
  "Samox",
  "datadrift-example",
  "036f9d6b685ee02a14faa70ed05e0bd60650c477",
];

function App() {
  const [dualTableProps, setTableProps] = useState<DualTableProps | null>(null);

  useEffect(() => {
    const fetchCommitData = async () => {
      try {
        const files = await getCommitFiles(owner, repo, commitSHA);
        if (!files) {
          throw new Error("No files found");
        }
        if (files[0] && files[0].patch) {
          const { oldData, newData } = parsePatch(files[0].patch);
          console.log("oldData", oldData);
          console.log("newData", newData);
          setTableProps({ tableProps1: oldData, tableProps2: newData });
        }
      } catch (error) {
        console.error("Error fetching GitHub commit data:", error);
      }
    };

    fetchCommitData().catch(console.error);
  }, []);

  return (
    <>
      <a href={`https://github.com/${owner}/${repo}/commit/${commitSHA}`}>
        Link to commit {`${owner}/${repo}/commit/${commitSHA}`}
      </a>
      {dualTableProps && <DualTable {...dualTableProps} />}
    </>
  );
}

export default App;
