import { Params, useLoaderData } from "react-router";
import { getCommitList } from "../services/data-drift";
import { CommitListItem } from "../components/Commits/CommitListItem";

function assertParamsIsDefined(
  params: Params<"installationId" | "owner" | "repo">
): asserts params is { installationId: string; owner: string; repo: string } {
  if (
    typeof params.installationId === "string" &&
    params.installationId.trim() !== "" &&
    typeof params.owner === "string" &&
    params.owner.trim() !== "" &&
    typeof params.repo === "string" &&
    params.repo.trim() !== ""
  ) {
    return;
  }
  throw new Error("Params is not defined");
}

const loader = async ({
  params,
}: {
  params: Params<"installationId" | "owner" | "repo">;
}) => {
  assertParamsIsDefined(params);
  console.log(params);
  const result = await getCommitList(params);
  return { data: result.data, params };
};

type LoaderData = Awaited<ReturnType<typeof loader>>;

const ddCommitDiffUrlFactory = (params: {
  installationId: string;
  owner: string;
  repo: string;
  commitSha: string;
}) => {
  return `/report/${params.installationId}/${params.owner}/${params.repo}/commit/${params.commitSha}`;
};

const DriftListPage = () => {
  const { data, params } = useLoaderData() as LoaderData;
  return (
    <div>
      {data.map((commit) => {
        const isDrift = commit.commit.message.includes("Drift");
        const commitUrl = ddCommitDiffUrlFactory({
          ...params,
          commitSha: commit.sha,
        });
        return (
          <CommitListItem
            type={isDrift ? "Drift" : "New Data"}
            date={
              commit.commit.author?.date
                ? new Date(commit.commit.author?.date)
                : null
            }
            name={commit.commit.message}
            commitUrl={commitUrl}
          />
        );
      })}
    </div>
  );
};

DriftListPage.loader = loader;

export default DriftListPage;
