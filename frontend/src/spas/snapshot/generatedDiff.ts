export type GeneratedDiff = {
  unique_key: {
    [key: number]: string;
  };
  record_status: {
    [key: number]: "before" | "after";
  };
  [key: string]: {
    [key: number]: string;
  };
};

export const sample_generated_diff = {
  unique_key: {
    0: "38c5f2ff-df60-46a7-b4c5-c5b0fee1d96f",
    1: "c427b416-f07c-495b-a1e5-b7f4b8e4a1f9",
    2: "38c5f2ff-df60-46a7-b4c5-c5b0fee1d96f",
    3: "800b9fa2-831c-4bae-8b09-c2f77ab1c07b",
    4: "c427b416-f07c-495b-a1e5-b7f4b8e4a1f9",
  },
  booking_date: {
    0: 1698796800000,
    1: 1699488000000,
    2: 1698796800000,
    3: 1669593600000,
    4: 1699488000000,
  },
  metric_value: { 0: 8.93, 1: 4.64, 2: 7.93, 3: 2.76, 4: 5.64 },
  created_at: {
    0: 1701085940365,
    1: 1701085940365,
    2: 1701085940365,
    3: 1701182689444,
    4: 1701085940365,
  },
  updated_at: {
    0: 1701085940365,
    1: 1701085940365,
    2: 1701182689444,
    3: 1701182689444,
    4: 1701182689444,
  },
  dbt_scd_id: {
    0: "76639c1a3d0aa1b0ff092a0e2487d9ff",
    1: "942c80ba457aea0a7db3c2e166be45b9",
    2: "6c0b4d87821c0cd01096878e8ae37034",
    3: "d76e429bc368364a4d83887f6b49471f",
    4: "c10d65d0e20a00986cfe9d68aa0a422c",
  },
  dbt_updated_at: {
    0: 1701085940365,
    1: 1701085940365,
    2: 1701182689444,
    3: 1701182689444,
    4: 1701182689444,
  },
  dbt_valid_from: {
    0: 1701085940365,
    1: 1701085940365,
    2: 1701182689444,
    3: 1701182689444,
    4: 1701182689444,
  },
  dbt_valid_to: {
    0: 1701182689444,
    1: 1701182689444,
    2: null,
    3: null,
    4: null,
  },
  record_status: {
    0: "before",
    1: "before",
    2: "after",
    3: "after",
    4: "after",
  },
} as const;
