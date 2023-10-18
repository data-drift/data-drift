import { Params, useLoaderData } from "react-router-dom";
import { getTable } from "../services/data-drift";
import { CommitListItem } from "../components/Commits/CommitListItem";

const loader = async ({ params }: { params: Params<string> }) => {
  const tableName = params.tableName as string;
  const tableInfo = await getTable(tableName);
  return tableInfo;
};

type LoaderData = Awaited<ReturnType<typeof loader>>;

const TablePage = () => {
  const loader = useLoaderData() as LoaderData;
  return (
    <div style={{ textAlign: "justify" }}>
      <h1>Store: {loader.data.store}</h1>
      <h2>Table: {loader.data.table}</h2>
      <h3>Columns:</h3>
      <ul>
        {loader.data.tableColumns.map((table) => (
          <li key={table} style={{ textAlign: "justify" }}>
            <a href={`./${loader.data.table}/metrics/${table}`}>{table}</a>
          </li>
        ))}
      </ul>
      <h3>History:</h3>
      {loader.data.commits.length > 0 ? (
        loader.data.commits.map((commit) => {
          const isDrift = commit.Message.includes("DRIFT");
          const commitUrl = "";
          return (
            <CommitListItem
              key={commit.Sha}
              type={isDrift ? "Drift" : "New Data"}
              date={new Date(commit.Date)}
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
  );
};

TablePage.loader = loader;

export default TablePage;
