import { Params, useLoaderData } from "react-router-dom";
import { getTableComparisonFromApi } from "../services/data-drift";
import { parsePatch } from "../services/patch.mapper";
import { DiffTable } from "./DisplayCommit/DiffTable";

const CompareCommits = () => {
  const loaderData = useLoaderData() as LoaderData;
  return (
    <>
      <p>
        Comparing data from{" "}
        <span title={loaderData.fromDate.toString()}>
          {loaderData.fromDate.toLocaleDateString()}
        </span>
        {" to "}
        <span title={loaderData.toDate.toString()}>
          {loaderData.toDate.toLocaleDateString()}
        </span>
      </p>
      <DiffTable dualTableProps={loaderData.dualTableProps} />
    </>
  );
};

const loader = async ({
  params: { installationId, owner, repo },
}: {
  params: Params<"installationId" | "owner" | "repo">;
}) => {
  const searchParams = new URLSearchParams(window.location.search);
  const beginDate = searchParams.get("begin-date");
  const endDate = searchParams.get("end-date");
  const table = searchParams.get("table");
  if (!beginDate || !endDate || !table || !installationId || !owner || !repo)
    throw new Error("Missing params");
  const comparison = await getTableComparisonFromApi({
    owner,
    repo,
    table,
    beginDate,
    endDate,
  });

  const { oldData, newData } = parsePatch(
    comparison.data.patch,
    comparison.data.headers
  );
  const dualTableProps = {
    tableProps1: oldData,
    tableProps2: newData,
  };
  return {
    dualTableProps,
    fromDate: new Date(comparison.data.baseCommitDateISO8601),
    toDate: new Date(comparison.data.headCommitDateISO8601),
  };
};

type LoaderData = Awaited<ReturnType<typeof loader>>;

CompareCommits.loader = loader;

export default CompareCommits;
