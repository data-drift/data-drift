import { Params, useLoaderData } from "react-router";
import { DDConfig, configQuery, getCommitList } from "../services/data-drift";
import { CommitList } from "../components/Commits/CommitList";
import DriftCard from "../components/Commits/DriftCard";
import styled from "@emotion/styled";
import { QueryClient } from "@tanstack/react-query";

function assertParamsIsDefined(
  params: Params<"owner" | "repo">
): asserts params is { installationId: string; owner: string; repo: string } {
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

function extractParentsFromConfig(
  config: DDConfig,
  filepath: string
): NonNullable<DDConfig["metrics"][number]["upstreamFiles"]> {
  const metric = config.metrics.find((metric) => metric.filepath === filepath);
  return metric ? metric.upstreamFiles || [] : [];
}

function queryParamsAreDefined(params: Record<string, string>): params is {
  periodKey: string;
  filepath: string;
  driftDate: string;
} {
  return "periodKey" in params && "filepath" in params && "driftDate" in params;
}

const loader =
  (queryClient: QueryClient) =>
  async ({ params }: { params: Params<"owner" | "repo"> }) => {
    assertParamsIsDefined(params);
    const result = await getCommitList(params);
    const query = configQuery(params);
    const maybeConfig = queryClient.getQueryData<DDConfig>(query.queryKey);
    const config =
      maybeConfig !== undefined
        ? maybeConfig
        : await queryClient.fetchQuery(query);
    const urlParams = Object.fromEntries(
      new URLSearchParams(window.location.search)
    );
    if (queryParamsAreDefined(urlParams)) {
      const urlParamsWithParent = {
        ...urlParams,
        parentData: extractParentsFromConfig(config, urlParams.filepath),
      };
      return {
        data: result.data,
        params,
        urlParams: urlParamsWithParent,
      };
    }
    return {
      data: result.data,
      params,
      urlParams: {
        periodKey: "",
        filepath: "",
        driftDate: "",
        parentData: [],
      },
    };
  };

type LoaderData = Awaited<ReturnType<ReturnType<typeof loader>>>;

const DriftListContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 16px;
  width: "100%";
`;

const DriftListPage = () => {
  const { data, params, urlParams } = useLoaderData() as LoaderData;
  const driftCardState = DriftCard.useState({ ...urlParams });
  return (
    <DriftListContainer>
      {urlParams && urlParams.filepath.length > 1 && (
        <DriftCard {...driftCardState} />
      )}
      <CommitList data={data} params={params} filters={driftCardState} />
    </DriftListContainer>
  );
};

DriftListPage.loader = loader;

export default DriftListPage;
