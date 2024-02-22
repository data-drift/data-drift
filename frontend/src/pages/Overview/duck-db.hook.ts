import * as duckdb from "@duckdb/duckdb-wasm";
import duckdb_wasm from "@duckdb/duckdb-wasm/dist/duckdb-mvp.wasm?url";
import mvp_worker from "@duckdb/duckdb-wasm/dist/duckdb-browser-mvp.worker.js?url";
import duckdb_wasm_eh from "@duckdb/duckdb-wasm/dist/duckdb-eh.wasm?url";
import eh_worker from "@duckdb/duckdb-wasm/dist/duckdb-browser-eh.worker.js?url";
import React, { useEffect, useRef, useState } from "react";
import * as arrow from "apache-arrow";
import { DualTableProps } from "../../components/Table/DualTable";

const MANUAL_BUNDLES = {
  mvp: {
    mainModule: duckdb_wasm,
    mainWorker: mvp_worker,
  },
  eh: {
    mainModule: duckdb_wasm_eh,
    mainWorker: eh_worker,
  },
} as const;

let singletonDb: duckdb.AsyncDuckDBConnection | null = null;

const useDuckDB = () => {
  const [db, setDb] = useState<duckdb.AsyncDuckDBConnection | null>(null);

  useEffect(() => {
    const initDuckDB = async () => {
      if (singletonDb) {
        setDb(singletonDb);
        return;
      }
      try {
        const bundle = await duckdb.selectBundle(MANUAL_BUNDLES);
        if (!bundle.mainWorker) {
          throw new Error("No worker found in the selected bundle");
        }
        const worker = new Worker(bundle.mainWorker);
        const logger = new duckdb.ConsoleLogger();
        const dbInstance = new duckdb.AsyncDuckDB(logger, worker);
        await dbInstance.instantiate(bundle.mainModule, bundle.pthreadWorker);
        const connection = await dbInstance.connect();
        singletonDb = connection;

        await connection.query(
          `CREATE TABLE people(id INTEGER, name VARCHAR);`
        );
        await connection.query(`INSERT INTO people VALUES (1, 'Mark');`);
        await connection.query(`INSERT INTO people VALUES (2, 'Phil');`);
        await connection.query(`INSERT INTO people VALUES (3, 'Roger');`);

        setDb(connection);
      } catch (error) {
        console.error("Failed to initialize DuckDB:", error);
      }
    };

    void initDuckDB();

    return () => {
      // TODO implement Clean up. Close db ?
    };
  }, []);

  return { db, useDbQuery };
};

const useDbQuery = <
  T extends { [key: string]: arrow.DataType<arrow.Type, any> }
>(
  sql: string,
  db: duckdb.AsyncDuckDBConnection | null
) => {
  const [result, setResult] = useState<arrow.Table<T> | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const hasTableBeenLoaded = useRef<boolean>(false);

  useEffect(() => {
    const queryAndSetResult = async () => {
      if (db) {
        try {
          const queryResult = await db.query<T>(sql);
          setResult(queryResult);
          setLoading(false);
        } catch (error) {
          setLoading(false);
          setError(error as Error);
          console.error("Query error:", error);
        }
      } else {
        setLoading(true);
      }
    };
    void queryAndSetResult();
  }, [sql, db]);

  return { result, loading, error, hasTableBeenLoaded };
};

export const mapQueryResultToPeople = (queryResult: arrow.Table) => {
  const queryResultToArray = queryResult.toArray();
  const result = queryResultToArray.map((arrow) => {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-unsafe-call
    const row = arrow.toArray() as [arrow.Int32, arrow.Utf8];
    return {
      id: row[0],
      name: row[1],
    };
  });
  return result;
};

export const useLoadSnapshotData = (
  dualTableData: DualTableProps | undefined,
  db: duckdb.AsyncDuckDBConnection | null,
  hasTableBeenLoaded: React.MutableRefObject<boolean>
) => {
  const hasEffectRun = hasTableBeenLoaded;

  useEffect(() => {
    const handleDualTableLoaded = async () => {
      if (dualTableData && db && !hasEffectRun.current) {
        console.log("dualTableData", dualTableData);
        await loadSnapshotData(dualTableData, db);
        hasEffectRun.current = true;
      }
    };

    void handleDualTableLoaded();
  }, [dualTableData, db, hasEffectRun]);
};

export const loadSnapshotData = async (
  dualTableProps: DualTableProps,
  db: duckdb.AsyncDuckDBConnection
) => {
  const oldData = dualTableProps.tableProps1;
  const createOldTableSql = `CREATE TABLE old_snapshot (${oldData.headers.join(
    " VARCHAR, "
  )} VARCHAR)`;
  await db.query(createOldTableSql);
  for (const row of oldData.data) {
    if (row.isEmpty || row.isEllipsis) {
      continue;
    }
    const insertSql = `INSERT INTO old_snapshot VALUES ('${row.data.join(
      "', '"
    )}')`;
    await db.query(insertSql);
  }

  const newData = dualTableProps.tableProps1;

  const createNewTableSql = `CREATE TABLE new_snapshot (${newData.headers.join(
    " VARCHAR, "
  )} VARCHAR)`;
  await db.query(createNewTableSql);
  for (const row of newData.data) {
    if (row.isEmpty || row.isEllipsis) {
      continue;
    }
    const insertSql = `INSERT INTO old_snapshot VALUES ('${row.data.join(
      "', '"
    )}')`;
    await db.query(insertSql);
  }
  return;
};

export default useDuckDB;
