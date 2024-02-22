import CodeEditor from "@uiw/react-textarea-code-editor";
import { useState } from "react";
import { useLoadSnapshotData } from "./duck-db.hook";
import { DualTableProps } from "../../components/Table/DualTable";
import * as duckdb from "@duckdb/duckdb-wasm";
import DuckDbProvider from "../../components/DuckDb/DuckDbProvider";

type SqlEditorProps = {
  dualTable: DualTableProps;
  db: duckdb.AsyncDuckDBConnection;
};

const SqlEditor = ({ dualTable }: SqlEditorProps) => {
  const db = DuckDbProvider.useDuckDb();
  useLoadSnapshotData(dualTable, db);

  const [sql, setSQL] = useState("");
  const [isRunning, setIsRunning] = useState(false);

  const onValidation = async () => {
    if (!db) {
      console.error("DuckDB is not initialized.");
      return;
    }
    try {
      setIsRunning(true);
      const results = await db.query(sql);
      console.log(results);
      setIsRunning(false);
    } catch (error) {
      console.error(error);
      setIsRunning(false);
    }
  };
  return (
    <div>
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
        }}
      />
      <button onClick={onValidation} disabled={isRunning}>
        {isRunning ? "Running..." : "Run"}
      </button>
    </div>
  );
};

export default SqlEditor;
