import { getCommitFiles, getCsvHeaders } from "../services/github";
import { DualTable, DualTableProps } from "../components/Table/DualTable";
import { parsePatch } from "../services/patch.mapper";
import { useLoaderData } from "react-router";

export interface CommitInfo {
  owner: string;
  repo: string;
  commitSHA: string;
}

const getOldAndNewDataFromGithub = async ({
  params: { owner, repo, commitSHA },
}: {
  params: { owner: string; repo: string; commitSHA: string };
}) => {
  const files = await getCommitFiles(owner, repo, commitSHA);
  if (!files) {
    throw new Error("No files found");
  }
  const file = files[0];
  if (file && file.patch) {
    const headers = await getCsvHeaders(file.contents_url);
    const { oldData, newData } = parsePatch(file.patch, headers);
    return { tableProps1: oldData, tableProps2: newData };
  }
};

function DisplayCommit() {
  const results = useLoaderData() as DualTableProps;
  console.log("results", results);

  return <>{results && <DualTable {...results} />}</>;
}

DisplayCommit.githubLoader = getOldAndNewDataFromGithub;

export default DisplayCommit;
