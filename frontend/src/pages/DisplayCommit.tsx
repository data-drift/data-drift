import { useEffect, useMemo, useState } from "react";
import { getCommitFiles, getCsvHeaders } from "../services/github";
import { DualTable, DualTableProps } from "../components/Table/DualTable";
import { parsePatch } from "../services/patch.mapper";

export interface CommitInfo {
  owner: string;
  repo: string;
  commitSHA: string;
}

function DisplayCommit() {
  const [dualTableProps, setTableProps] = useState<DualTableProps | null>(null);
  const [commitInfo, setCommitInfo] = useState<CommitInfo | null>(null);
  const pathArray = useMemo(() => window.location.pathname.split("/"), []);

  useEffect(() => {
    const fetchCommitData = async () => {
      try {
        const owner = pathArray[1];
        const repo = pathArray[2];
        const commitSHA = pathArray[4];
        if (!owner || !repo || !commitSHA) {
          return;
        }
        setCommitInfo({ owner, repo, commitSHA });
        const files = await getCommitFiles(owner, repo, commitSHA);
        if (!files) {
          throw new Error("No files found");
        }
        const file = files[0];
        if (file && file.patch) {
          const headers = await getCsvHeaders(file.contents_url);
          const { oldData, newData } = parsePatch(file.patch, headers);
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

  return <>{dualTableProps && <DualTable {...dualTableProps} />}</>;
}

export default DisplayCommit;
