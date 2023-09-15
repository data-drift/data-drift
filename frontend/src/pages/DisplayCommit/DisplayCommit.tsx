import { getCommitFiles, getCsvHeaders } from "../../services/github";
import { parsePatch } from "../../services/patch.mapper";
import { Params, useLoaderData } from "react-router";
import { getPatchAndHeader } from "../../services/data-drift";
import styled from "@emotion/styled";
import { DiffTable } from "./DiffTable";

export interface CommitParam {
  owner: string;
  repo: string;
  commitSHA: string;
}

const StyledButton = styled.button`
  padding: 4px 16px;
  color: ${(props) => props.theme.colors.text};
  border-radius: 0px;
  background-color: ${(props) => props.theme.colors.background2};
  border: 1px solid ${(props) => props.theme.colors.text};
`;

function assertParamsIsCommitInfo(params: Params<string>): CommitParam {
  const { owner, repo, commitSHA } = params;
  if (!owner || !repo || !commitSHA) {
    throw new Error("Invalid params");
  }
  return { owner, repo, commitSHA };
}

function assertParamsHasInstallationIs(
  params: Params<string>
): CommitParam & { installationId: string } {
  const { installationId, owner, repo, commitSHA } = params;
  if (!installationId || !owner || !repo || !commitSHA) {
    throw new Error("Invalid params");
  }
  return { installationId, owner, repo, commitSHA };
}

const getCommitDiffFromGithub = async ({
  params,
}: {
  params: Params<string>;
}) => {
  const { owner, repo, commitSHA } = assertParamsIsCommitInfo(params);
  const files = await getCommitFiles(owner, repo, commitSHA);
  if (!files) {
    throw new Error("No files found");
  }
  const file = files[0];
  if (file && file.patch) {
    const headers = await getCsvHeaders(file.contents_url);
    const { oldData, newData } = parsePatch(file.patch, headers);
    return {
      data: { tableProps1: oldData, tableProps2: newData },
      params: { owner, repo, commitSHA },
    };
  }
};

const getCommitDiffFromDataDrift = async ({
  params,
}: {
  params: Params<string>;
}) => {
  const { installationId, owner, repo, commitSHA } =
    assertParamsHasInstallationIs(params);

  const { patch, headers, ...commitInfo } = await getPatchAndHeader({
    installationId,
    owner,
    repo,
    commitSHA,
  });

  const { oldData, newData } = parsePatch(patch, headers);
  return {
    data: { tableProps1: oldData, tableProps2: newData, commitInfo },
    params: { owner, repo, commitSHA, installationId },
  };
};

type LoaderData = Awaited<
  ReturnType<typeof getCommitDiffFromGithub | typeof getCommitDiffFromDataDrift>
>;

const StyledSpan = styled.span`
  padding: 8px;
  align-self: flex-start;
  display: flex;
  gap: 8px;
  align-items: center;
`;

const StyledIcon = styled.img`
  filter: invert(1);
  height: 24px;
  vertical-align: middle;
`;

const ddCommitListUrlFactory = (params: {
  installationId: string;
  owner: string;
  repo: string;
}) => {
  return `/report/${params.installationId}/${params.owner}/${params.repo}/commits`;
};

function DisplayCommit() {
  const results = useLoaderData() as LoaderData;

  return (
    <>
      {results && "commitInfo" in results.data && (
        <StyledSpan>
          <b>{results.data.commitInfo.filename}</b> -{" "}
          <b>{results.data.commitInfo.date.toLocaleDateString()}</b>
          <a href={results.data.commitInfo.commitLink}>
            <StyledIcon src="/src/assets/github-mark.svg" alt="GitHub" />
          </a>
          {"installationId" in results.params && (
            <a href={ddCommitListUrlFactory(results.params)}>
              <StyledButton>View list of commits</StyledButton>
            </a>
          )}
        </StyledSpan>
      )}
      {results && results.data && <DiffTable dualTableProps={results.data} />}
    </>
  );
}

DisplayCommit.githubLoader = getCommitDiffFromGithub;
DisplayCommit.dataDriftLoader = getCommitDiffFromDataDrift;

export default DisplayCommit;
