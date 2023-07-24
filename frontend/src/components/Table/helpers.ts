import { Datum } from "./Table";

export const mapStringArrayToDatum = (data: string[]): Datum[] => {
  return data.map((value, index) => ({ value, isEmphasized: index === 1 }));
};
