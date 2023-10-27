import { Params, useLoaderData } from "react-router-dom";
import {
  DDConfig,
  getCommitList,
  getCommitListLocalStrategy,
  getConfig,
} from "../../services/data-drift";

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
  assertParamsIsDefined(params);

  const config = await getConfig(params);

  return {
    params,
    config,
    strategy: "github",
  };
};

export const localStrategyLoader = async ({
  params,
}: {
  params: Params<"tableName">;
}) => {
  const tableName = params.tableName || "";

  const config: DDConfig = {
    metrics: [
      {
        metricName: tableName,
        filepath: tableName,
        upstreamFiles: [],
        dateColumnName: "",
        KPIColumnName: "",
        timeGrains: [],
        dimensions: [],
      },
    ],
  };

  return {
    params,
    config,
    strategy: "local",
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
