import { Params, useLoaderData } from "react-router-dom";
import { getTableComparisonFromApi } from "../services/data-drift";
import { parsePatch } from "../services/patch.mapper";
import { DiffTable } from "./DisplayCommit/DiffTable";
import { DualTableProps } from "../components/Table/DualTable";

const CompareCommits = () => {
  const loaderData = useLoaderData() as DualTableProps;
  return <DiffTable dualTableProps={loaderData} />;
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
    installationId,
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
  return dualTableProps;
};

CompareCommits.loader = loader;

export default CompareCommits;
