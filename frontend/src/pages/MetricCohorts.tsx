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
  owner?: string;
  repo?: string;
  installationId?: string;
  metricName: string;
  timegrain: Timegrain;
} {
  const { installationId, owner, repo, metricName, timegrain } = params;
  if (!metricName || !timegrain) {
    throw new Error("Invalid params");
  }
  if (!installationId && (!owner || !repo)) {
    throw new Error(
      "Either installationId or both owner and repo must be defined"
    );
  }
  assertTimegrain(timegrain);

  return { installationId, owner, repo, metricName, timegrain };
}

const ScrollableContainer = styled.div`
  width: 100%;
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
