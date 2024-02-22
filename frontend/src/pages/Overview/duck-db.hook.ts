import * as duckdb from "@duckdb/duckdb-wasm";
import duckdb_wasm from "@duckdb/duckdb-wasm/dist/duckdb-mvp.wasm?url";
import mvp_worker from "@duckdb/duckdb-wasm/dist/duckdb-browser-mvp.worker.js?url";
import duckdb_wasm_eh from "@duckdb/duckdb-wasm/dist/duckdb-eh.wasm?url";
import eh_worker from "@duckdb/duckdb-wasm/dist/duckdb-browser-eh.worker.js?url";
import { useEffect, useState } from "react";
import * as arrow from "apache-arrow";

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

type People = {
  id: arrow.Int32;
  name: arrow.Utf8;
};

const useDuckDB = () => {
  const [db, setDb] = useState<duckdb.AsyncDuckDBConnection | null>(null);

  useEffect(() => {
    const initDuckDB = async () => {
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
      // Clean up
      //   db?.terminate(); // Assuming terminate() is the correct method to clean up. Adjust as necessary.
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

  return { result, loading, error };
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

export default useDuckDB;
