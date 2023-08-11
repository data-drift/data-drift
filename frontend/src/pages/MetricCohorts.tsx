import { Params, useLoaderData } from "react-router";
import { StepChart, StepChartProps } from "../components/Charts/StepChart";
import {
  Timegrain,
  assertTimegrain,
  getMetricCohorts,
} from "../services/data-drift";

const getMetricCohortsData = ({
  params,
}: {
  params: Params<string>;
}): StepChartProps => {
  const typedParams = assertParamsHasNeededProperties(params);
  const result = getMetricCohorts(typedParams);
  console.log(result);
  return {
    metricNames: ["2022-01"],
    data: [
      { daysSinceFirstReport: 12, "2022-01": 12 },
      { daysSinceFirstReport: 13, "2022-01": 12 },
    ],
  };
};

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
