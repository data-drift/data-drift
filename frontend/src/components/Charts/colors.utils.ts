import { scaleLinear } from "d3-scale";
import { YearMonthString } from "./StepChart";

const colorSelector = (year: string) => {
  switch (year) {
    case "2022":
      return ["red", "blue"];
    case "2023":
      return ["green", "yellow"];
    default:
      return ["black", "white"];
  }
};

export const getMetricColor = (yearMonthString: YearMonthString) => {
  const [year, month] = yearMonthString.split("-");
  console.log(year, month);
  const scale = scaleLinear([0, 11], colorSelector(year));

  console.log(scale(parseInt(month, 10)));
  return scale(parseInt(month, 10));
};
