import { Params, useLoaderData } from "react-router-dom";
import {
  DDConfig,
  configQuery,
  getCommitList,
  getCommitListLocalStrategy,
  getMeasurement,
  getPatchAndHeader,
} from "../../services/data-drift";
import { QueryClient } from "@tanstack/react-query";
import { parsePatch } from "../../services/patch.mapper";

enum Strategy {
  Github = "github",
  Local = "local",
}
function assertParamsIsDefined(
  params: Params<"owner" | "repo">
): asserts params is { owner: string; repo: string } {
  if (
    typeof params.owner === "string" &&
    params.owner.trim() !== "" &&
    typeof params.repo === "string" &&
    params.repo.trim() !== ""
  ) {
    return;
  }
  throw new Error("Params is not defined");
}

export const loader =
  (queryClient: QueryClient) =>
  async ({ params }: { params: Params<"owner" | "repo"> }) => {
    assertParamsIsDefined(params);
    const query = configQuery(params);
    const maybeConfig = queryClient.getQueryData<DDConfig>(query.queryKey);
    const config =
      maybeConfig !== undefined
        ? maybeConfig
        : await queryClient.fetchQuery(query);

    return {
      params,
      config,
      strategy: Strategy.Github,
    } as const;
  };

export const localStrategyLoader = ({
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
    params: { ...params, tableName },
    config,
    strategy: Strategy.Local,
  } as const;
};

type LoaderData = Awaited<
  ReturnType<ReturnType<typeof loader> | typeof localStrategyLoader>
>;

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

const fetchCommits = async (
  strategy:
    | {
        strategy: Strategy.Local;
        params: { tableName: string };
      }
    | {
        strategy: Strategy.Github;
        params: { owner: string; repo: string };
      },
  currentDate: Date
): Promise<
  {
    commit: { message: string; author: { date?: string } | null };
    sha: string;
  }[]
> => {
  switch (strategy.strategy) {
    case "local": {
      const result = await getCommitListLocalStrategy(
        strategy.params.tableName,
        currentDate.toISOString().substring(0, 10)
      );

      const mappedCommits = result.data.Measurements.map((commit) => ({
        commit: {
          message: commit.Message,
          author: {
            date: commit.Date,
          },
        },
        sha: commit.Sha,
      })) satisfies {
        commit: { message: string; author: { date?: string } | null };
        sha: string;
      }[];

      return mappedCommits;
    }
    case "github": {
      const result = await getCommitList(
        strategy.params,
        currentDate.toISOString().substring(0, 10)
      );
      return result.data;
    }
    default:
      throw new Error("Strategy not supported");
  }
};

export const fetchCommitsQuery = (
  strategy:
    | {
        strategy: Strategy.Local;
        params: { tableName: string };
      }
    | {
        strategy: Strategy.Github;
        params: { owner: string; repo: string };
      },
  currentDate: Date
) => ({
  queryKey: ["commit", strategy, currentDate],
  queryFn: () => fetchCommits(strategy, currentDate),
});

const fetchCommitPatch = async (
  strategy:
    | {
        strategy: Strategy.Local;
        params: { tableName: string };
      }
    | {
        strategy: Strategy.Github;
        params: { owner: string; repo: string };
      },
  selectedCommit: string
) => {
  switch (strategy.strategy) {
    case Strategy.Local: {
      const measurementResults = await getMeasurement(
        "default",
        strategy.params.tableName,
        selectedCommit
      );
      const { oldData, newData } = parsePatch(
        measurementResults.data.Patch,
        measurementResults.data.Headers
      );
      const dualTableProps = {
        tableProps1: oldData,
        tableProps2: newData,
      };
      return dualTableProps;
    }
    case Strategy.Github: {
      const patchAndHeader = await getPatchAndHeader({
        owner: strategy.params.owner,
        repo: strategy.params.repo,
        commitSHA: selectedCommit,
      });
      const { oldData, newData } = parsePatch(
        patchAndHeader.patch,
        patchAndHeader.headers
      );
      const dualTableProps = {
        tableProps1: oldData,
        tableProps2: newData,
      };
      return dualTableProps;
    }
    default: {
      const unhandeldStrategy: never = strategy;
      console.error("Strategy not supported", unhandeldStrategy);
      throw new Error("Strategy not supported");
    }
  }
};

export const fetchCommitPatchQuery = (
  strategy:
    | {
        strategy: Strategy.Local;
        params: { tableName: string };
      }
    | {
        strategy: Strategy.Github;
        params: { owner: string; repo: string };
      },
  selectedCommit?: string | null
) => ({
  queryKey: ["commit", "patch", strategy, selectedCommit],
  queryFn: () => fetchCommitPatch(strategy, selectedCommit as string),
  enabled: !!selectedCommit,
});
