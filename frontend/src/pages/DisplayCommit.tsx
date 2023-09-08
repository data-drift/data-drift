import { getCommitFiles, getCsvHeaders } from "../services/github";
import { DualTable } from "../components/Table/DualTable";
import { parsePatch } from "../services/patch.mapper";
import { Params, useLoaderData } from "react-router";
import { getPatchAndHeader } from "../services/data-drift";
import styled from "@emotion/styled";

export interface CommitParam {
  owner: string;
  repo: string;
  commitSHA: string;
}

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
    return { data: { tableProps1: oldData, tableProps2: newData }, params };
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
    params,
  };
};

type LoaderData = Awaited<
  ReturnType<typeof getCommitDiffFromGithub | typeof getCommitDiffFromDataDrift>
>;

const StyledSpan = styled.span`
  padding: 8px;
`;

const ddCommitListUrlFactory = (params: {
  installationId: string;
  owner: string;
  repo: string;
}) => {
  return `/report/${params.installationId}/${params.owner}/${params.repo}/commits}`;
};

function DisplayCommit() {
  const results = useLoaderData() as LoaderData;

  return (
    <>
      {results && "commitInfo" in results.data && (
        <StyledSpan>
          <b>
            {results.data.commitInfo.filename}{" "}
            {results.data.commitInfo.date.toLocaleDateString()}{" "}
          </b>
          <a href={results.data.commitInfo.commitLink}>github</a>
        </StyledSpan>
      )}
      {results && <DualTable {...results.data} />}
    </>
  );
}

DisplayCommit.githubLoader = getCommitDiffFromGithub;
DisplayCommit.dataDriftLoader = getCommitDiffFromDataDrift;

export default DisplayCommit;
