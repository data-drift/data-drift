import * as duckdb from "@duckdb/duckdb-wasm";
import { useEffect, useRef } from "react";
import * as arrow from "apache-arrow";
import { DualTableProps } from "../../components/Table/DualTable";

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
  db: duckdb.AsyncDuckDBConnection | null
) => {
  const hasEffectRun = useRef<boolean>(false);

  useEffect(() => {
    const handleDualTableLoaded = async () => {
      if (dualTableData && db && !hasEffectRun.current) {
        hasEffectRun.current = true;
        await loadSnapshotData(dualTableData, db);
      }
    };

    void handleDualTableLoaded();
  }, [dualTableData, db]);
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
  console.log("old_snapshot table created");
  for (const row of oldData.data) {
    if (row.isEmpty || row.isEllipsis) {
      continue;
    }
    const insertSql = `INSERT INTO old_snapshot VALUES ('${row.data
      .map((datum) => datum.value)
      .join("', '")}')`;
    await db.query(insertSql);
  }
  console.log("old_snapshot data inserted created");

  const newData = dualTableProps.tableProps2;
  const createNewTableSql = `CREATE TABLE new_snapshot (${newData.headers.join(
    " VARCHAR, "
  )} VARCHAR)`;
  await db.query(createNewTableSql);
  console.log("new_snapshot table created");
  for (const row of newData.data) {
    if (row.isEmpty || row.isEllipsis) {
      continue;
    }
    const insertSql = `INSERT INTO new_snapshot VALUES ('${row.data
      .map((datum) => datum.value)
      .join("', '")}')`;
    await db.query(insertSql);
  }
  console.log("new_snapshot data inserted created");
  return;
};
