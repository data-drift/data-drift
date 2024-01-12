import { parsePatch } from "../../services/patch.mapper";
import { Params, useLoaderData, defer } from "react-router";
import { getPatchAndHeader } from "../../services/data-drift";
import styled from "@emotion/styled";
import { DiffTable } from "./DiffTable";
import { toast } from "react-toastify";
import Loader from "../../components/Common/Loader";
import React from "react";
import { Await } from "react-router-dom";

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

function assertParamsIsCommitParam(params: Params<string>): CommitParam {
  const { owner, repo, commitSHA } = params;
  if (!owner || !repo || !commitSHA) {
    throw new Error("Invalid params");
  }
  return { owner, repo, commitSHA };
}

const getPatchFromApi = async ({ owner, repo, commitSHA }: CommitParam) => {
  const [{ patch, headers, patchToLarge, ...commitInfo }] = await Promise.all([
    getPatchAndHeader({
      owner,
      repo,
      commitSHA,
    }),
  ]);

  if (patchToLarge) {
    toast(
      "Diff is too large to display. Only showing partial diff. Display may be broken.",
      { autoClose: false }
    );
  }

  const { oldData, newData } = parsePatch(patch, headers);
  const data = { tableProps1: oldData, tableProps2: newData, commitInfo };
  return data;
};

const getCommitDiffFromDataDrift = ({ params }: { params: Params<string> }) => {
  const { owner, repo, commitSHA } = assertParamsIsCommitParam(params);

  const data = getPatchFromApi({ owner, repo, commitSHA });

  return defer({
    data: data,
    params: { owner, repo, commitSHA },
  });
};

type PatchFromApi = Awaited<ReturnType<typeof getPatchFromApi>>;

type LoaderData = ReturnType<typeof getCommitDiffFromDataDrift> & {
  params: CommitParam;
};

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

const ddCommitListUrlFactory = (
  params: {
    owner: string;
    repo: string;
  },
  queryParams?: { periodKey: string; filepath: string; driftDate: string }
) => {
  const url = `/report/${params.owner}/${params.repo}/commits`;
  if (queryParams) {
    const urlQueryParams = new URLSearchParams(queryParams).toString();
    return url + "?" + urlQueryParams;
  }
  return url;
};

const PageContainer = styled.div`
  height: 100vh;
  width: 100vw;
`;

function DisplayCommit() {
  const results = useLoaderData() as LoaderData;
  const resultsParam = results.params;

  const searchParams = new URLSearchParams(window.location.search);
  const periodKey = searchParams.get("periodKey") as string;

  return (
    <PageContainer>
      <React.Suspense fallback={<Loader />}>
        <Await resolve={results.data}>
          {(resultsData: PatchFromApi) => (
            <>
              <StyledSpan>
                <b>{resultsData.commitInfo.filename}</b> -{" "}
                <b>{resultsData.commitInfo.date.toLocaleDateString()}</b>
                <a href={resultsData.commitInfo.commitLink}>
                  <StyledIcon src="/github-mark.svg" alt="GitHub" />
                </a>
                <a
                  href={ddCommitListUrlFactory(resultsParam, {
                    periodKey,
                    filepath: resultsData.commitInfo.filename,
                    driftDate: resultsData.commitInfo.date.toISOString(),
                  })}
                >
                  <StyledButton>View list of commits</StyledButton>
                </a>
              </StyledSpan>
              <DiffTable dualTableProps={resultsData} />
            </>
          )}
        </Await>
      </React.Suspense>
    </PageContainer>
  );
}

DisplayCommit.dataDriftLoader = getCommitDiffFromDataDrift;

export default DisplayCommit;
