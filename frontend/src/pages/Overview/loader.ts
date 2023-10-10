import { Params, useLoaderData } from "react-router-dom";
import { getCommitList, getConfig } from "../../services/data-drift";

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

export const loader = async ({
  params,
}: {
  params: Params<"installationId" | "owner" | "repo">;
}) => {
  const searchParams = new URLSearchParams(window.location.search);
  const snapshotDate = searchParams.get("snapshotDate") || undefined;

  assertParamsIsDefined(params);
  const [result, config] = await Promise.all([
    getCommitList(params, snapshotDate),
    getConfig(params),
  ]);

  return {
    data: result.data,
    params,
    config,
  };
};

type LoaderData = Awaited<ReturnType<typeof loader>>;

function assertLoaderDataIsDefined(
  loaderData: unknown
): asserts loaderData is LoaderData {
  if (typeof loaderData === "object" && loaderData !== null) {
    return;
  }
  throw new Error("Loader data is not defined");
}

export const useOverviewLoaderData = () => {
  const loaderData = useLoaderData();
  assertLoaderDataIsDefined(loaderData);
  return loaderData;
};
