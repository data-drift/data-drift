import { Row, TableProps } from "../components/Table/Table";

export const parsePatch = (patch: string) => {
  const lines = patch.split("\n");
  const headersLine = lines.shift();
  if (!headersLine) throw new Error("No headers line found");
  const headerData = headersLine.match(
    /^@@ -(\d+),(\d+) \+(\d+),(\d+) @@ (.*)$/
  );
  let headerString: string;
  if (!headerData) {
    headerString = lines.shift() || "";
  } else {
    headerString = headerData[5];
  }

  if (!headerString) throw new Error("No header string found");
  const headers = headerString.split(",").map((header) => header.trim()) || [];

  const oldData: TableProps = { diffType: "removed", data: [], headers };
  const newData: TableProps = { diffType: "added", data: [], headers };

  const uniqueKeyAfter: Record<string, Row> = {};
  const uniqueKeyBefore: Record<string, Row> = {};

  let csvColumnsLength = 0;

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (csvColumnsLength === 0) {
      csvColumnsLength = line.split(",").length;
    }
    const lineData = line.substring(1);
    const uniqueKey = getUniqueKey(lineData);
    if (line.startsWith("-")) {
      uniqueKeyBefore[uniqueKey] = csvStringLineToRowData(lineData, true);
    } else if (line.startsWith("+")) {
      uniqueKeyAfter[uniqueKey] = csvStringLineToRowData(lineData, true);
    } else if (line.startsWith(" ")) {
      uniqueKeyBefore[uniqueKey] = csvStringLineToRowData(lineData);
      uniqueKeyAfter[uniqueKey] = csvStringLineToRowData(lineData);
    }
  }

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    const lineData = line.substring(1);
    const uniqueKey = getUniqueKey(lineData);
    if (line.startsWith("-")) {
      oldData.data.push(csvStringLineToRowData(lineData, true));
      if (!uniqueKeyAfter[uniqueKey]) {
        newData.data.push(emptyRow(csvColumnsLength));
      }
    } else if (line.startsWith("+")) {
      newData.data.push(csvStringLineToRowData(lineData, true));
      if (!uniqueKeyBefore[uniqueKey]) {
        oldData.data.push(emptyRow(csvColumnsLength));
      }
    } else if (line.startsWith(" ")) {
      oldData.data.push(csvStringLineToRowData(lineData));
      newData.data.push(csvStringLineToRowData(lineData));
    }
  }

  return { oldData, newData };
};

const getUniqueKey = (line: string) => {
  return line.split(",")[0];
};

const emptyRow = (csvColumnsLength: number): Row => ({
  data: Array.from({ length: csvColumnsLength }).map(() => ({ value: "_" })),
  isEmphasized: false,
});

const csvStringLineToRowData = (line: string, isEmphasized = false): Row => {
  return {
    data: line.split(",").map((value) => ({ value })),
    isEmphasized,
  };
};
