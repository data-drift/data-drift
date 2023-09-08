import { Params, useLoaderData } from "react-router";
import { getCommitList } from "../services/data-drift";

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
  return result.data;
};

type LoaderData = Awaited<ReturnType<typeof loader>>;

const DriftListPage = () => {
  const data = useLoaderData() as LoaderData;
  return (
    <div>
      Drift List Page
      {data.map((commit) => {
        return (
          <pre key={commit.sha}>{JSON.stringify(commit.commit.message)}</pre>
        );
      })}
    </div>
  );
};

DriftListPage.loader = loader;

export default DriftListPage;
