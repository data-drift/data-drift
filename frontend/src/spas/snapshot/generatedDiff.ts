import { DualTableProps } from "../../components/Table/DualTable";
import { Row } from "../../components/Table/Table";

export type SnapshotDiff = {
  unique_key: {
    [key: number]: string;
  };
  record_status: {
    [key: number]: "before" | "after";
  };
  dbt_scd_id: {
    [key: number]: string;
  };
  dbt_updated_at: {
    [key: number]: string;
  };
  dbt_valid_from: {
    [key: number]: string;
  };
  dbt_valid_to: {
    [key: number]: string | null;
  };
  [key: string]: {
    [key: number]: string | number | null;
  };
};

export const getHeaders = (diff: SnapshotDiff): string[] => {
  const headers = Object.keys(diff).filter(
    (key) =>
      key !== "record_status" &&
      key !== "dbt_scd_id" &&
      key !== "dbt_valid_from" &&
      key !== "dbt_valid_to" &&
      key !== "dbt_updated_at"
  );
  return headers;
};

export const mapSnapshotDiffToRows = (
  diff: SnapshotDiff
): { removedRows: Row[]; addedRows: Row[] } => {
  return { removedRows: [], addedRows: [] };
};

export const mapSnapshotDiffToDualTableProps = (
  diff: SnapshotDiff
): DualTableProps => {
  const headers = getHeaders(diff);
  const { removedRows, addedRows } = mapSnapshotDiffToRows(diff);
  return {
    tableProps1: {
      headers,
      data: removedRows,
      diffType: "removed",
    },
    tableProps2: {
      headers,
      data: addedRows,
      diffType: "added",
    },
  };
};

export const sampleSnapshotDiff = {
  unique_key: {
    "0": "38c5f2ff-df60-46a7-b4c5-c5b0fee1d96f",
    "1": "c427b416-f07c-495b-a1e5-b7f4b8e4a1f9",
    "2": "38c5f2ff-df60-46a7-b4c5-c5b0fee1d96f",
    "3": "800b9fa2-831c-4bae-8b09-c2f77ab1c07b",
    "4": "c427b416-f07c-495b-a1e5-b7f4b8e4a1f9",
  },
  booking_date: {
    "0": "2023-11-01T00:00:00.000",
    "1": "2023-11-09T00:00:00.000",
    "2": "2023-11-01T00:00:00.000",
    "3": "2022-11-28T00:00:00.000",
    "4": "2023-11-09T00:00:00.000",
  },
  metric_value: { "0": 8.93, "1": 4.64, "2": 7.93, "3": 2.76, "4": 5.64 },
  country_code: { "0": "BN", "1": "KN", "2": "BN", "3": "US", "4": "KN" },
  created_at: {
    "0": "2023-11-27T11:52:20.365",
    "1": "2023-11-27T11:52:20.365",
    "2": "2023-11-27T11:52:20.365",
    "3": "2023-11-28T14:44:49.444",
    "4": "2023-11-27T11:52:20.365",
  },
  updated_at: {
    "0": "2023-11-27T11:52:20.365",
    "1": "2023-11-27T11:52:20.365",
    "2": "2023-11-28T14:44:49.444",
    "3": "2023-11-28T14:44:49.444",
    "4": "2023-11-28T14:44:49.444",
  },
  dbt_scd_id: {
    "0": "76639c1a3d0aa1b0ff092a0e2487d9ff",
    "1": "942c80ba457aea0a7db3c2e166be45b9",
    "2": "6c0b4d87821c0cd01096878e8ae37034",
    "3": "d76e429bc368364a4d83887f6b49471f",
    "4": "c10d65d0e20a00986cfe9d68aa0a422c",
  },
  dbt_updated_at: {
    "0": "2023-11-27T11:52:20.365",
    "1": "2023-11-27T11:52:20.365",
    "2": "2023-11-28T14:44:49.444",
    "3": "2023-11-28T14:44:49.444",
    "4": "2023-11-28T14:44:49.444",
  },
  dbt_valid_from: {
    "0": "2023-11-27T11:52:20.365",
    "1": "2023-11-27T11:52:20.365",
    "2": "2023-11-28T14:44:49.444",
    "3": "2023-11-28T14:44:49.444",
    "4": "2023-11-28T14:44:49.444",
  },
  dbt_valid_to: {
    "0": "2023-11-28T14:44:49.444",
    "1": "2023-11-28T14:44:49.444",
    "2": null,
    "3": null,
    "4": null,
  },
  record_status: {
    "0": "before",
    "1": "before",
    "2": "after",
    "3": "after",
    "4": "after",
  },
} as const satisfies SnapshotDiff;
