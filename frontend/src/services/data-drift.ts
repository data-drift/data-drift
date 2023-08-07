import axios from "axios";
import { CommitInfo } from "../pages/DisplayCommit";

export const getPatchAndHeader = async (
  params: CommitInfo & { installationId: string }
) => {
  console.log("commitParams", params);
  const result = await axios.get<{ patch: string; headers: string[] }>(
    `https://data-drift.herokuapp.com/gh/${params.owner}/${params.repo}/commit/${params.commitSHA}`,
    { headers: { "Installation-Id": params.installationId } }
  );
  console.log("result", result);
  return {
    patch: result.data.patch,
    headers: result.data.headers.join(","),
  };
};
