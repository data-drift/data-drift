import { Params, useLoaderData } from "react-router";
import { getCommitList, getConfig } from "../services/data-drift";
import { CommitList } from "../components/Commits/CommitList";
import { DriftCard } from "../components/Commits/DriftCard";
import styled from "@emotion/styled";

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

function queryParamsAreDefined(params: Record<string, string>): params is {
  periodKey: string;
  fileName: string;
  driftDate: string;
} {
  return "periodKey" in params && "fileName" in params && "driftDate" in params;
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
  const urlParams = Object.fromEntries(
    new URLSearchParams(window.location.search)
  );
  if (queryParamsAreDefined(urlParams)) {
    const urlParamsWithParent = {
      ...urlParams,
      parentData: ["metrics/ride_daily_revenue.csv"],
      filepath: urlParams.fileName,
    };
    return {
      data: result.data,
      params,
      config,
      urlParams: urlParamsWithParent,
    };
  }
  return { data: result.data, params, config };
};

type LoaderData = Awaited<ReturnType<typeof loader>>;

const DriftListContainer = styled.div`
  display: flex;
  flex-direction: row;
  gap: 8px;
`;

const DriftListPage = () => {
  const { data, params, urlParams } = useLoaderData() as LoaderData;
  return (
    <DriftListContainer>
      {urlParams && (
        <DriftCard {...urlParams} parentData={urlParams.parentData} />
      )}
      <CommitList data={data} params={params} />
    </DriftListContainer>
  );
};

DriftListPage.loader = loader;

export default DriftListPage;
