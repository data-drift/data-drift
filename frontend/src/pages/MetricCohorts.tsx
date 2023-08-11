import { Params, useLoaderData } from "react-router";
import { StepChart, StepChartProps } from "../components/Charts/StepChart";

const getMetricCohortsData = ({
  params,
}: {
  params: Params<string>;
}): StepChartProps => {
  const typedParams = assertParamsHasNeededProperties(params);
  console.log(typedParams);
  return {
    metricNames: ["2022-01"],
    data: [
      { daysSinceFirstReport: 12, "2022-01": 12 },
      { daysSinceFirstReport: 13, "2022-01": 12 },
    ],
  };
};

// Define the custom type
type Timegrain = "quarter" | "month" | "week" | "day";

// The assertion function
function assertTimegrain(value: string): asserts value is Timegrain {
  if (
    value !== "quarter" &&
    value !== "month" &&
    value !== "week" &&
    value !== "day"
  ) {
    throw new Error("Value is not a valid time unit!");
  }
}

function assertParamsHasNeededProperties(params: Params<string>): {
  installationId: string;
  metricName: string;
  timegrain: Timegrain;
} {
  const { installationId, metricName, timegrain } = params;
  if (!installationId || !metricName || !timegrain) {
    throw new Error("Invalid params");
  }
  assertTimegrain(timegrain);

  return { installationId, metricName, timegrain };
}

const MetricCohort = () => {
  const { metricNames, data } = useLoaderData() as StepChartProps;
  return <StepChart metricNames={metricNames} data={data} />;
};

MetricCohort.loader = getMetricCohortsData;
export default MetricCohort;
