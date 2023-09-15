import { parsePatch } from "./patch.mapper";
import * as addedAndRemoved from "./patch.mapper.tests.use-cases/added-and-removed.ts";

describe("parsePatch", () => {
  it("should detect added and removed lines with isEmphasized on the rows", () => {
    const patch = addedAndRemoved.patch;
    const { oldData, newData } = parsePatch(patch, [
      "unique_key",
      "name",
      "date",
      "age",
    ]);

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
    const { oldData, newData } = parsePatch(patch, [
      "unique_key",
      "name",
      "date",
      "age",
    ]);

    expect(oldData.headers).toEqual(["unique_key", "name", "date", "age"]);
    expect(newData.headers).toEqual(["unique_key", "name", "date", "age"]);
  });

  it("should parse the headers in the first line", () => {
    const patch =
      "@@ -2,9 +2,9 @@\n unique_key,name,date,age\n 2022-12-Alice,Alice,2022-12,25";
    const { oldData, newData } = parsePatch(patch, [
      "unique_key",
      "name",
      "date",
      "age",
    ]);

    expect(oldData.headers).toEqual(["unique_key", "name", "date", "age"]);
    expect(newData.headers).toEqual(["unique_key", "name", "date", "age"]);
  });

  it("should detect modified lines with isEmphasized on the rows", () => {
    const patch =
      "@@ -2,9 +2,9 @@ unique_key,name,date,age\n-2022-12-Alice,Alice,2022-12,25\n+2022-12-Alice,Alice,2022-12,26";
    const { oldData, newData } = parsePatch(patch, [
      "unique_key",
      "name",
      "date",
      "age",
    ]);
    expect(oldData.data[0].data[3]).toStrictEqual({
      value: "25",
      type: "number",
      isEmphasized: true,
    });
    expect(newData.data[0].data[3]).toStrictEqual({
      value: "26",
      type: "number",
      isEmphasized: true,
    });
  });

  it("should add ellipsis between hunks", () => {
    const patch =
      "@@ -11537,1 +11537,1 @@ unique_key,date,metric_value,country_code,category\n 1d6b2882-9a26-41d6-a976-dff75dfdd5e2,2004-05-11,3.93,PT,Category C\n@@ -23809,1 +23809,1 @@ unique_key,date,metric_value,country_code,category\n 3d16e1f2-a2fc-4b6a-a00b-fbc4e16aa485,2012-05-20,2.17,CZ,Category C";
    const { oldData, newData } = parsePatch(patch, [
      "unique_key",
      "date",
      "metric_value",
      "country_code",
      "category",
    ]);
    expect(oldData.data[1].isEllipsis).toStrictEqual(true);
    expect(newData.data[1].isEllipsis).toStrictEqual(true);
  });

  it("should handle large datasets", () => {
    const patch =
      "@@ -11537,7 +11537,7 @@ unique_key,date,metric_value,country_code,category\n 1d6b2882-9a26-41d6-a976-dff75dfdd5e2,2004-05-11,3.93,PT,Category C\n 1d6b2f8c-13c7-4c33-9835-d7e9e5448f81,2012-12-16,2.18,AF,Category C\n 1d6b412b-e1be-4566-b27f-d6d2591be73e,2019-11-01,4.0,SI,Category A\n-1d6c2c23-39cb-418a-9e09-e22b8ca3601f,1993-10-25,2.59,PE,Category B\n+1d6c2c23-39cb-418a-9e09-e22b8ca3601f,1993-10-25,7.22,PE,Category B\n 1d6c6e9c-3f2b-4f36-af8f-4ece80e1c456,2021-11-07,1.72,ME,Category C\n 1d6cd9a0-4f30-452f-8dda-34777e2814e0,2002-01-04,7.89,AM,Category C\n 1d6e65d0-0b9b-4b76-9517-ed65dd038b14,2009-01-25,4.63,GW,Category C\n@@ -23809,7 +23809,7 @@ unique_key,date,metric_value,country_code,category\n 3d1634ea-2648-41cf-8bc4-02063d57d6d4,2015-05-19,8.43,KE,Category A\n 3d16de2d-c626-46b7-90bb-211be6e38f31,2004-02-27,0.69,JO,Category A\n 3d16e1f2-a2fc-4b6a-a00b-fbc4e16aa485,2012-05-20,2.17,CZ,Category C";

    const results = parsePatch(patch, [
      "unique_key",
      "date",
      "metric_value",
      "country_code",
      "category",
    ]);
    expect(results).toMatchSnapshot();
  });

  it("should handle when headers change", () => {
    const patch =
      "@@ -1,31 +1,31 @@\n-unique_key,date_month,country,mrr_bop,resurrected_mrr,downsell_mrr,fail_to_engage_mrr,churn_mrr,mrr_eop\n-2022-01-01__FR,2022-01-01,FR,11.11,11.0,-11.0,-11.0,-11.0,11.11\n+unique_key,date_month,country,mrr_bop,downsell_mrr,fail_to_engage_mrr,churn_mrr,mrr_eop\n+2022-01-01__FR,2022-01-01,FR,22.22,-22.0,-22.0,-22.0,22.22";

    const results = parsePatch(
      patch,
      "unique_key,date_month,country,mrr_bop,resurrected_mrr,downsell_mrr,fail_to_engage_mrr,churn_mrr,mrr_eop".split(
        ","
      )
    );
    expect(results).toMatchSnapshot();
  });

  it("should handle when headers reverse order", () => {
    const patch =
      "@@ -1,31 +1,31 @@\n-unique_key,date_month,country,mrr_bop\n-2022-01-01__FR,2022-01-01,FR,11.11\n+unique_key,date_month,mrr_bop,country\n+2022-01-01__FR,2022-01-01,22.22,FR";

    const results = parsePatch(patch, [
      "unique_key",
      "date_month",
      "mrr_bop",
      "country",
    ]);
    expect(results).toMatchSnapshot();
  });

  it("should handle when headers are in patch", () => {
    const patch =
      "@@ -1,31 +1,31 @@\n unique_key,date_month,country,mrr_bop\n-2022-01-01__FR,2022-01-01,FR,11.11\n unique_key,date_month,country,mrr_bop\n+2022-01-01__FR,2022-01-01,FR,22.22";

    const results = parsePatch(patch, [
      "unique_key",
      "date_month",
      "mrr_bop",
      "country",
    ]);
    expect(results).toMatchSnapshot();
  });
});
