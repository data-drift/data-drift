import { Params, useLoaderData } from "react-router-dom";
import { getTable } from "../services/data-drift";
import { CommitListItem } from "../components/Commits/CommitListItem";
import StarUs from "../components/Common/StarUs";

const loader = async ({ params }: { params: Params<string> }) => {
  const tableName = params.tableName as string;
  const tableInfo = await getTable(tableName);
  return tableInfo;
};

type LoaderData = Awaited<ReturnType<typeof loader>>;

const TablePage = () => {
  const loader = useLoaderData() as LoaderData;
  return (
    <div
      style={{
        textAlign: "justify",
        width: "100%",
        boxSizing: "border-box",
        padding: "0 24px",
      }}
    >
      <h1
        style={{
          width: "100%",
          display: "flex",
          justifyContent: "space-between",
        }}
      >
        <span>Store: {loader.data.store}</span>
        <StarUs />
      </h1>
      <h2>Table: {loader.data.table}</h2>
      <h3>Columns:</h3>
      <ul>
        {loader.data.tableColumns.map((metric) => (
          <li key={metric} style={{ textAlign: "justify" }}>
            <a href={`./${loader.data.table}/metrics/${metric}`}>{metric}</a>
          </li>
        ))}
      </ul>
      <h3>History:</h3>
      <div style={{ maxWidth: "fit-content" }}>
        {loader.data.commits.length > 0 ? (
          loader.data.commits.map((commit) => {
            const isDrift = commit.Message.includes("DRIFT");
            const commitDate = new Date(commit.Date);
            const commitUrl = `/tables/${
              loader.data.table
            }/history?snapshotDate=${
              commitDate.toISOString().split("T")[0]
            }&commitSha=${commit.Sha}`;
            return (
              <CommitListItem
                key={commit.Sha}
                type={isDrift ? "Drift" : "New Data"}
                date={commitDate}
                name={commit.Message}
                commitUrl={commitUrl}
                isParentData={false}
              />
            );
          })
        ) : (
          <div
            style={{
              border: "1px solid #ccc",
              padding: "16px",
              borderRadius: "0",
              marginBottom: "16px",
              display: "flex",
              flexDirection: "column",
              alignItems: "flex-start",
            }}
          >
            No commits found
          </div>
        )}
      </div>
    </div>
  );
};

TablePage.loader = loader;

export default TablePage;
