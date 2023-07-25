import { useEffect, useMemo, useState } from "react";
import "./App.css";
import { getCommitFiles } from "./services/github";
import { DualTable, DualTableProps } from "./components/Table/DualTable";
import { parsePatch } from "./services/patch.mapper";

interface CommitInfo {
  owner: string;
  repo: string;
  commitSHA: string;
}

function App() {
  const [dualTableProps, setTableProps] = useState<DualTableProps | null>(null);
  const [commitInfo, setCommitInfo] = useState<CommitInfo | null>(null);
  const pathArray = useMemo(() => window.location.pathname.split("/"), []);

  useEffect(() => {
    const fetchCommitData = async () => {
      try {
        const owner = pathArray[1];
        const repo = pathArray[2];
        const commitSHA = pathArray[4];
        setCommitInfo({ owner, repo, commitSHA });
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
  }, [pathArray]);

  return commitInfo ? (
    <>
      <a
        href={`https://github.com/${commitInfo.owner}/${commitInfo.repo}/commit/${commitInfo.commitSHA}`}
      >
        Link to commit{" "}
      </a>
      {dualTableProps && <DualTable {...dualTableProps} />}
    </>
  ) : (
    <div>Could not parse URL with format /$owner/$repo/commit/$commitSHA</div>
  );
}

export default App;
