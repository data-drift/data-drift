import { Row, TableProps } from "../components/Table/Table";

export const parsePatch = (patch: string, headers: string[]) => {
  if (patch === "") {
    throw new Error("Diff is too large to be displayed");
  }
  let oldHeaders = headers;
  const lines = patch.split("\n");
  const headersLine = lines.shift();
  if (!headersLine) throw new Error("No headers line found");
  let firstAddedLineShouldBeSkiped = false;
  const headerData = headersLine.match(
    /^@@ -(\d+),(\d+) \+(\d+),(\d+) @@ (.*)$/
  );
  if (!headerData) {
    const oldHeadersStringWithModifier = lines.shift();
    const modifier = oldHeadersStringWithModifier?.substring(0, 1);
    if (modifier === "-") {
      firstAddedLineShouldBeSkiped = true;
    }
    const oldHeadersString = oldHeadersStringWithModifier?.substring(1);
    oldHeaders =
      oldHeadersString?.split(",").map((header) => header.trim()) || [];
  }

  const oldData: TableProps = {
    diffType: "removed",
    data: [],
    headers: oldHeaders,
  };
  const newData: TableProps = { diffType: "added", data: [], headers };

  const rowByUniqueKeyAfter: Record<string, Row> = {};
  const rowByUniqueKeyBefore: Record<string, Row> = {};

  let csvColumnsLength = 0;

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (csvColumnsLength === 0) {
      csvColumnsLength = line.split(",").length;
    }
    const lineData = line.substring(1);
    const uniqueKey = getUniqueKey(lineData);
    if (line.startsWith("-")) {
      rowByUniqueKeyBefore[uniqueKey] = csvStringLineToRowData(lineData, true);
    } else if (line.startsWith("+")) {
      rowByUniqueKeyAfter[uniqueKey] = csvStringLineToRowData(lineData, true);
    } else if (line.startsWith(" ")) {
      rowByUniqueKeyBefore[uniqueKey] = csvStringLineToRowData(lineData);
      rowByUniqueKeyAfter[uniqueKey] = csvStringLineToRowData(lineData);
    }
  }

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    const lineData = line.substring(1);
    const uniqueKey = getUniqueKey(lineData);
    if (line.startsWith("-")) {
      oldData.data.push(csvStringLineToRowData(lineData, true));
      if (!rowByUniqueKeyAfter[uniqueKey]) {
        newData.data.push(emptyRow(csvColumnsLength));
      } else {
        const emphasizedCellIndexes = getCellIndexesToEmphasize(
          csvStringLineToRowData(lineData, true),
          rowByUniqueKeyAfter[uniqueKey],
          getNewIndexFromOldIndex(oldHeaders, headers)
        );
        emphasizedCellIndexes.forEach((index) => {
          oldData.data[oldData.data.length - 1].data[index].isEmphasized = true;
        });
      }
    } else if (line.startsWith("+")) {
      if (firstAddedLineShouldBeSkiped) {
        firstAddedLineShouldBeSkiped = false;
        continue;
      }
      newData.data.push(csvStringLineToRowData(lineData, true));
      if (!rowByUniqueKeyBefore[uniqueKey]) {
        oldData.data.push(emptyRow(csvColumnsLength));
      } else {
        const emphasizedCellIndexes = getCellIndexesToEmphasize(
          csvStringLineToRowData(lineData, true),
          rowByUniqueKeyBefore[uniqueKey],
          getOldIndexFromNewIndex(oldHeaders, headers)
        );
        emphasizedCellIndexes.forEach((index) => {
          if (newData.data[newData.data.length - 1].data[index])
            newData.data[newData.data.length - 1].data[index].isEmphasized =
              true;
        });
      }
    } else if (line.startsWith(" ")) {
      oldData.data.push(csvStringLineToRowData(lineData));
      newData.data.push(csvStringLineToRowData(lineData));
    } else if (line.startsWith("@@")) {
      oldData.data.push({
        isEllipsis: true,
        data: [],
      });
      newData.data.push({
        isEllipsis: true,
        data: [],
      });
    }
  }

  newData.data.forEach((row, index) => {
    row.data.forEach((cell, cellIndex) => {
      if (cell.isEmphasized && cell.type === "number") {
        if (
          !oldData.data[index].data[
            getOldIndexFromNewIndex(oldHeaders, headers)(cellIndex)
          ]
        ) {
          return;
        }

        cell.diffValue =
          Number(cell.value) -
          Number(
            oldData.data[index].data[
              getOldIndexFromNewIndex(oldHeaders, headers)(cellIndex)
            ].value
          );
        oldData.data[index].data[
          getOldIndexFromNewIndex(oldHeaders, headers)(cellIndex)
        ].diffValue = -cell.diffValue;
      }
    });
  });

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
  let fields: string[] = [];

  if (!line.includes('"')) {
    fields = line.split(",");
  } else {
    let insideQuotes = false;
    let fieldStart = 0;

    for (let i = 0; i <= line.length; i++) {
      if (i === line.length || (line[i] === "," && !insideQuotes)) {
        fields.push(line.substring(fieldStart, i));
        fieldStart = i + 1;
      } else if (line[i] === '"') {
        insideQuotes = !insideQuotes;
      }
    }
  }

  return {
    data: fields.map((value) => ({
      value:
        value.startsWith('"') && value.endsWith('"')
          ? value.slice(1, -1)
          : value,
      type: Number.isNaN(Number(value)) ? "string" : "number",
    })),
    isEmphasized,
  };
};

const getCellIndexesToEmphasize = (
  row: Row,
  rowToCompare: Row,
  indexMapper: (number: number) => number | undefined
): number[] => {
  const differentIndexes: number[] = [];
  row.data.forEach((cell, index) => {
    const indexToCompare = indexMapper(index);
    if (
      indexToCompare &&
      cell.value !== rowToCompare.data[indexToCompare]?.value
    ) {
      differentIndexes.push(index);
    }
  });
  return differentIndexes;
};

const getNewIndexFromOldIndex = (
  oldHeaders: string[],
  newHeaders: string[]
) => {
  return (oldIndex: number) => {
    const oldHeader = oldHeaders[oldIndex];
    const newIndex = newHeaders.indexOf(oldHeader);
    return newIndex;
  };
};

const getOldIndexFromNewIndex = (
  oldHeaders: string[],
  newHeaders: string[]
) => {
  return (newIndex: number) => {
    const newHeader = newHeaders[newIndex];
    const oldIndex = oldHeaders.indexOf(newHeader);
    return oldIndex;
  };
};
