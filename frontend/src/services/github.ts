import axios from "axios";
import { Endpoints } from "@octokit/types";

type CommitResponse =
  Endpoints["GET /repos/{owner}/{repo}/commits/{ref}"]["response"];

export const getCommitFiles = async (
  owner: string,
  repo: string,
  commitSHA: string
) => {
  const response = await axios.get<CommitResponse["data"]>(
    `https://api.github.com/repos/${owner}/${repo}/commits/${commitSHA}`
  );
  if (response.status !== 200) {
    throw new Error("Error fetching commit content");
  }
  return response?.data?.files;
};
