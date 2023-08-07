export const getPatchAndHeader = (commitParams: any) => {
  console.log("commitParams", commitParams);
  return {
    patch:
      "@@ -2,9 +2,9 @@ unique_key,name,date,age\n 2022-12-Alice,Alice,2022-12,25",
    headers: "unique_key,name,date,age",
  };
};
