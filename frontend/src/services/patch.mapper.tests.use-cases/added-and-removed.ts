import { Row } from "../../components/Table/Table";

export const patch =
  "@@ -2,9 +2,9 @@ unique_key,name,date,age\n 2022-12-Alice,Alice,2022-12,25\n 2023-01-Bob,Bob,2023-01,30\n 2023-01-Charlie,Charlie,2023-01,36\n+2023-01-Charline,Charline,2023-01,42\n 2023-02-Antoine,Antoine,2023-02,40\n-2023-02-Didier,Didier,2023-02,40\n+2023-02-Didier,Didier,2023-02,43\n 2023-02-Philipe,Philipe,2023-02,42\n-2023-03-Clement,Clement,2023-03,45\n 2023-03-Cyril,Cyril,2023-03,45\n 2023-03-Victor,Victor,2023-03,46";

export const expectedOldData: Row[] = [
  {
    isEmphasized: false,
    data: [
      { value: "2022-12-Alice", type: "string" },
      { value: "Alice", type: "string" },
      { value: "2022-12", type: "string" },
      { value: "25", type: "number" },
    ],
  },
  {
    isEmphasized: false,
    data: [
      { value: "2023-01-Bob", type: "string" },
      { value: "Bob", type: "string" },
      { value: "2023-01", type: "string" },
      { value: "30", type: "number" },
    ],
  },
  {
    isEmphasized: false,
    data: [
      { value: "2023-01-Charlie", type: "string" },
      { value: "Charlie", type: "string" },
      { value: "2023-01", type: "string" },
      { value: "36", type: "number" },
    ],
  },
  {
    isEmphasized: false,
    data: [{ value: "_" }, { value: "_" }, { value: "_" }, { value: "_" }],
  }, // empty line here
  {
    isEmphasized: false,
    data: [
      { value: "2023-02-Antoine", type: "string" },
      { value: "Antoine", type: "string" },
      { value: "2023-02", type: "string" },
      { value: "40", type: "number" },
    ],
  },
  {
    isEmphasized: true,
    data: [
      { value: "2023-02-Didier", type: "string" },
      { value: "Didier", type: "string" },
      { value: "2023-02", type: "string" },
      { value: "40", type: "number", isEmphasized: true },
    ],
  },
  {
    isEmphasized: false,
    data: [
      { value: "2023-02-Philipe", type: "string" },
      { value: "Philipe", type: "string" },
      { value: "2023-02", type: "string" },
      { value: "42", type: "number" },
    ],
  },
  {
    isEmphasized: true,
    data: [
      { value: "2023-03-Clement", type: "string" },
      { value: "Clement", type: "string" },
      { value: "2023-03", type: "string" },
      { value: "45", type: "number" },
    ],
  },
  {
    isEmphasized: false,
    data: [
      { value: "2023-03-Cyril", type: "string" },
      { value: "Cyril", type: "string" },
      { value: "2023-03", type: "string" },
      { value: "45", type: "number" },
    ],
  },
  {
    isEmphasized: false,
    data: [
      { value: "2023-03-Victor", type: "string" },
      { value: "Victor", type: "string" },
      { value: "2023-03", type: "string" },
      { value: "46", type: "number" },
    ],
  },
];

export const expectedNewData: Row[] = [
  {
    isEmphasized: false,
    data: [
      { value: "2022-12-Alice", type: "string" },
      { value: "Alice", type: "string" },
      { value: "2022-12", type: "string" },
      { value: "25", type: "number" },
    ],
  },
  {
    isEmphasized: false,
    data: [
      { value: "2023-01-Bob", type: "string" },
      { value: "Bob", type: "string" },
      { value: "2023-01", type: "string" },
      { value: "30", type: "number" },
    ],
  },
  {
    isEmphasized: false,
    data: [
      { value: "2023-01-Charlie", type: "string" },
      { value: "Charlie", type: "string" },
      { value: "2023-01", type: "string" },
      { value: "36", type: "number" },
    ],
  },
  {
    isEmphasized: true,
    data: [
      { value: "2023-01-Charline", type: "string" },
      { value: "Charline", type: "string" },
      { value: "2023-01", type: "string" },
      { value: "42", type: "number" },
    ],
  },
  {
    isEmphasized: false,
    data: [
      { value: "2023-02-Antoine", type: "string" },
      { value: "Antoine", type: "string" },
      { value: "2023-02", type: "string" },
      { value: "40", type: "number" },
    ],
  },
  {
    isEmphasized: true,
    data: [
      { value: "2023-02-Didier", type: "string" },
      { value: "Didier", type: "string" },
      { value: "2023-02", type: "string" },
      { value: "43", type: "number", isEmphasized: true },
    ],
  },
  {
    isEmphasized: false,
    data: [
      { value: "2023-02-Philipe", type: "string" },
      { value: "Philipe", type: "string" },
      { value: "2023-02", type: "string" },
      { value: "42", type: "number" },
    ],
  },
  {
    isEmphasized: false,
    data: [{ value: "_" }, { value: "_" }, { value: "_" }, { value: "_" }],
  }, // empty line here
  {
    isEmphasized: false,
    data: [
      { value: "2023-03-Cyril", type: "string" },
      { value: "Cyril", type: "string" },
      { value: "2023-03", type: "string" },
      { value: "45", type: "number" },
    ],
  },
  {
    isEmphasized: false,
    data: [
      { value: "2023-03-Victor", type: "string" },
      { value: "Victor", type: "string" },
      { value: "2023-03", type: "string" },
      { value: "46", type: "number" },
    ],
  },
];
