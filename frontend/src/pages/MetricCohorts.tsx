import { useLoaderData } from "react-router";
import { StepChart, StepChartProps } from "../components/Charts/StepChart";

const getMetricCohortsData = (): StepChartProps => {
  return {
    metricNames: ["2022-01"],
    data: [{ daysSinceFirstReport: 12, "2022-01": 12 }],
  };
};

const MetricCohort = () => {
  const { metricNames, data } = useLoaderData() as StepChartProps;
  console.log(metricNames, data);
  return <StepChart metricNames={metricNames} data={data} />;
};

MetricCohort.loader = getMetricCohortsData;
export default MetricCohort;
