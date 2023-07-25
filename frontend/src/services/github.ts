import axios from "axios";
import { Endpoints } from "@octokit/types";

type CommitResponse =
  Endpoints["GET /repos/{owner}/{repo}/commits/{ref}"]["response"];

const getRequestHeaders = () => {
  const githubAccessToken = localStorage.getItem("github_token");
  if (githubAccessToken) {
    return {
      Authorization: `bearer ${githubAccessToken}`,
    };
  } else {
    return {};
  }
};

export const getCommitFiles = async (
  owner: string,
  repo: string,
  commitSHA: string
) => {
  const headers = getRequestHeaders();
  const response = await axios.get<CommitResponse["data"]>(
    `https://api.github.com/repos/${owner}/${repo}/commits/${commitSHA}`,
    { headers }
  );
  if (response.status !== 200) {
    throw new Error("Error fetching commit content");
  }
  return response?.data?.files;
};
