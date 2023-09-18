import { Params, useLoaderData } from "react-router";
import { getCommitList } from "../services/data-drift";
import { CommitList } from "../components/Commits/CommitList";

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
  const result = await getCommitList(params);
  return { data: result.data, params };
};

type LoaderData = Awaited<ReturnType<typeof loader>>;

const DriftListPage = () => {
  const { data, params } = useLoaderData() as LoaderData;
  return (
    <div>
      <CommitList data={data} params={params} />
    </div>
  );
};

DriftListPage.loader = loader;

export default DriftListPage;
