import { Params, useLoaderData } from "react-router";
import { StepChart, StepChartProps } from "../components/Charts/StepChart";
import {
  Timegrain,
  assertTimegrain,
  getMetricCohorts,
} from "../services/data-drift";
import { mapCohortsMetricsMetadataToStepChartProps } from "../services/data-drift.mappers";
import styled from "@emotion/styled";

const getMetricCohortsData = async ({
  params,
}: {
  params: Params<string>;
}): Promise<StepChartProps> => {
  const typedParams = assertParamsHasNeededProperties(params);
  const result = await getMetricCohorts(typedParams);
  const { metricNames, data } = mapCohortsMetricsMetadataToStepChartProps(
    result.data.cohortsMetricsMetadata
  );
  return { metricNames, data };
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

const ScrollableContainer = styled.div`
  overflow-x: scroll;
  height: 260px;
`;

const MetricCohort = () => {
  const { metricNames, data } = useLoaderData() as StepChartProps;
  return (
    <ScrollableContainer>
      <StepChart metricNames={metricNames} data={data} />
    </ScrollableContainer>
  );
};

MetricCohort.loader = getMetricCohortsData;
export default MetricCohort;
