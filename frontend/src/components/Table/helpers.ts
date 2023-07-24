import { Datum } from "./Table";

export const mapStringArrayToDatum = (data: string[]): Datum[] => {
  return data.map((value) => ({ value, isEmphasized: false }));
};
