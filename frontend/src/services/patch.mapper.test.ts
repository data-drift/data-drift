import { parsePatch } from "./patch.mapper";
import * as addedAndRemoved from "./patch.mapper.tests.use-cases/added-and-removed.ts";

describe("parsePatch", () => {
  it("should detect added and removed lines with isEmphasized on the rows", () => {
    const patch = addedAndRemoved.patch;
    const { oldData, newData } = parsePatch(patch);

    expect(oldData.diffType).toBe("removed");
    expect(oldData.headers).toEqual(["unique_key", "name", "date", "age"]);

    expect(oldData.data).toEqual(addedAndRemoved.expectedOldData);

    expect(newData.diffType).toBe("added");
    expect(newData.headers).toEqual(["unique_key", "name", "date", "age"]);
    expect(newData.data).toEqual(addedAndRemoved.expectedNewData);
  });

  it("should parse the headers in the patch header", () => {
    const patch =
      "@@ -2,9 +2,9 @@ unique_key,name,date,age\n 2022-12-Alice,Alice,2022-12,25";
    const { oldData, newData } = parsePatch(patch);

    expect(oldData.headers).toEqual(["unique_key", "name", "date", "age"]);
    expect(newData.headers).toEqual(["unique_key", "name", "date", "age"]);
  });

  it("should parse the headers in the first line", () => {
    const patch =
      "@@ -2,9 +2,9 @@\n unique_key,name,date,age\n 2022-12-Alice,Alice,2022-12,25";
    const { oldData, newData } = parsePatch(patch);

    expect(oldData.headers).toEqual(["unique_key", "name", "date", "age"]);
    expect(newData.headers).toEqual(["unique_key", "name", "date", "age"]);
  });

  it("should detect modified lines with isEmphasized on the rows", () => {
    const patch =
      "@@ -2,9 +2,9 @@ unique_key,name,date,age\n-2022-12-Alice,Alice,2022-12,25\n+2022-12-Alice,Alice,2022-12,26";
    const { oldData, newData } = parsePatch(patch);
    expect(oldData.data[0].data[3]).toStrictEqual({
      value: "25",
      isEmphasized: true,
    });
    expect(newData.data[0].data[3]).toStrictEqual({
      value: "26",
      isEmphasized: true,
    });
  });
});
