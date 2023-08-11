import axios from "axios";
import { CommitInfo } from "../pages/DisplayCommit";

const DATA_DRIFT_API_URL = "https://data-drift.herokuapp.com";

export const getPatchAndHeader = async (
  params: CommitInfo & { installationId: string }
) => {
  const result = await axios.get<{ patch: string; headers: string[] }>(
    `${DATA_DRIFT_API_URL}/gh/${params.owner}/${params.repo}/commit/${params.commitSHA}`,
    { headers: { "Installation-Id": params.installationId } }
  );
  return {
    patch: result.data.patch,
    headers: result.data.headers,
  };
};
