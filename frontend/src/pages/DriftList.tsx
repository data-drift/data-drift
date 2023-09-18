import { Params, useLoaderData } from "react-router";
import { getCommitList, getConfig } from "../services/data-drift";
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
  const [result, config] = await Promise.all([
    getCommitList(params),
    getConfig(params),
  ]);
  return { data: result.data, params, config };
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
