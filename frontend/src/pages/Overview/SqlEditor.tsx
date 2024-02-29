import CodeEditor from "@uiw/react-textarea-code-editor";
import { useState } from "react";
import { useLoadSnapshotData } from "./duck-db.hook";
import { DualTableProps } from "../../components/Table/DualTable";
import * as duckdb from "@duckdb/duckdb-wasm";
import DuckDbProvider from "../../components/DuckDb/DuckDbProvider";
import { sqlToDualTableMapper } from "./sql-to-dual-table.mapper";

type SqlEditorProps = {
  dualTable: DualTableProps;
  db: duckdb.AsyncDuckDBConnection;
  setQueryResult: (result: DualTableProps) => void;
};

const SqlEditor = ({ dualTable, setQueryResult }: SqlEditorProps) => {
  const db = DuckDbProvider.useDuckDb();
  useLoadSnapshotData(dualTable, db);

  const [sql, setSQL] = useState("SELECT * FROM snapshot;");
  const [isRunning, setIsRunning] = useState(false);
  const [queryError, setQueryError] = useState<string | null>(null);

  const onValidation = () => {
    void handleValidation();
  };

  const handleValidation = async () => {
    if (!db) {
      console.error("DuckDB is not initialized.");
      return;
    }
    try {
      setIsRunning(true);
      const oldSql = sql.replace("snapshot", "old_snapshot");
      const newSql = sql.replace("snapshot", "new_snapshot");
      const oldResults = await db.query(oldSql);
      const oldRows = {
        values: oldResults.toArray().map(Object.fromEntries) as (Record<
          string,
          string
        > & { unique_key: string })[],
        columns: oldResults.schema.fields.map((d) => d.name),
      };
      const newResults = await db.query(newSql);
      const newRows = {
        values: newResults.toArray().map(Object.fromEntries) as (Record<
          string,
          string
        > & { unique_key: string })[],
        columns: newResults.schema.fields.map((d) => d.name),
      };
      const uniqueKeysSet = new Set([
        ...oldRows.values.map((row) => row["unique_key"]),
        ...newRows.values.map((row) => row["unique_key"]),
      ]);
      const uniqueKeys = Array.from(uniqueKeysSet);
      const queryResultDualTable = sqlToDualTableMapper(
        uniqueKeys,
        oldRows,
        newRows,
        dualTable
      );
      console.log("dualTable", queryResultDualTable);
      setQueryResult(queryResultDualTable);
      setIsRunning(false);
      setQueryError(null);
    } catch (error: any) {
      console.error(error);
      setQueryError((error as { message: string }).message);
      setIsRunning(false);
    }
  };
  return (
    <div style={{ height: "100%", display: "flex" }}>
      <CodeEditor
        value={sql}
        language="sql"
        placeholder="Please enter SQL code."
        onChange={(evn) => setSQL(evn.target.value)}
        padding={15}
        style={{
          backgroundColor: "#f3f3f3",
          width: "100%",
          height: "100%",
          fontSize: "16px",
          fontFamily:
            "ui-monospace,SFMono-Regular,SF Mono,Consolas,Liberation Mono,Menlo,monospace",
          color: "black",
          flex: 2,
        }}
      />
      <div style={{ flex: 1, paddingLeft: "4px" }}>
        <button onClick={onValidation} disabled={isRunning}>
          {isRunning ? "Running..." : "Run"}
        </button>
        {queryError && <div style={{ color: "red" }}>{queryError}</div>}
      </div>
    </div>
  );
};

export default SqlEditor;
