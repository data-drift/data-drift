import styled from "@emotion/styled";
import TrendChip from "../components/Charts/TrendChip";
import DualMetricBarChart from "../components/Charts/DualMetricBarChart";
import { Tooltip } from "../components/Common/Tooltip";

const PageContainer = styled.div`
  height: 100vh;
  width: 100vw;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  padding: 0 ${({ theme }) => theme.spacing(6)};
  align-items: flex-start;
  background-color: ${({ theme }) => theme.colors.background};
`;

const Separator = styled.div`
  width: 100%;
  border-bottom: 1px solid ${({ theme }) => theme.colors.text};
  margin-bottom: ${({ theme }) => theme.spacing(2)};
`;

const DriftDetailContainer = styled.div`
  width: 100%;
  background-color: ${({ theme }) => theme.colors.background2};
  padding: ${({ theme }) => theme.spacing(2)};
  box-sizing: border-box;
  clip-path: ${({ theme }) => theme.upLeftClipping};
  display: flex;
  flex-direction: row;
  gap: ${({ theme }) => theme.spacing(20)};
  overflow: visible;
`;

const SubSectionContainer = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: flex-start;
  margin-bottom: ${({ theme }) => theme.spacing(2)};
`;

const BlackContainer = styled.div`
  width: 100%;
  background-color: black;
  padding: ${({ theme }) => theme.spacing(2)};
  box-sizing: border-box;
  margin-top: ${({ theme }) => theme.spacing(2)};
  text-align: start;
  flex-direction: row;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: ${({ theme }) => theme.spacing(1)};
`;

const DrillDownButton = styled.button`
  padding: 8px 16px;
  background-color: ${({ theme }) => theme.colors.primary};
  color: black;
  border-radius: 0px;
  font-family: monospace;
  margin-left: auto;
  margin-top: auto;
  display: flex;
  flex-direction: row;
`;

const TansparentContainer = styled.div`
  width: 100%;
  padding: ${({ theme }) => theme.spacing(2)} 0;
  box-sizing: border-box;
  margin-top: ${({ theme }) => theme.spacing(2)};
  text-align: start;
`;

const DualBarChartContainer = styled.div`
  width: 100%;
  overflow-x: auto;
`;

const data = [
  {
    name: "MRR Jan 2023",
    before: 100395.76,
    after: 101395.76,
    percentageChange: 1,
  },
  {
    name: "MRR Feb 2023",
    before: 101395.76,
    after: 100295.76,
    percentageChange: -1.08,
  },
  {
    name: "MRR Mar 2023",
    before: 100295.76,
    after: 100295.76,
    percentageChange: 0,
  },
  {
    name: "MRR Apr 2023",
    before: 100295.76,
    after: 101142.12,
    percentageChange: 0.84,
  },
  {
    name: "MRR May 2023",
    before: 101142.12,
    after: 100092.34,
    percentageChange: -1.04,
  },
  {
    name: "MRR Jun 2023",
    before: 100092.34,
    after: 100092.34,
    percentageChange: 0,
  },
  {
    name: "MRR Jul 2023",
    before: 100092.34,
    after: 101042.18,
    percentageChange: 0.95,
  },
  {
    name: "MRR Aug 2023",
    before: 101042.18,
    after: 100395.56,
    percentageChange: -0.64,
  },
  {
    name: "MRR Sep 2023",
    before: 100395.56,
    after: 100395.56,
    percentageChange: 0,
  },
  {
    name: "MRR Oct 2023",
    before: 100395.56,
    after: 101695.34,
    percentageChange: 1.29,
  },
];

const DriftOverviewPage = () => {
  return (
    <PageContainer>
      <h1>Drift: MRR monthly</h1>
      <Separator />
      <DriftDetailContainer>
        <SubSectionContainer>
          <strong>Metric Metadata</strong>
          <BlackContainer>
            <Tooltip>
              Stabilization time: 10d
              <span>
                The monthly metric stops moving 10 days after the end of the
                month
              </span>
            </Tooltip>
          </BlackContainer>
          <BlackContainer>
            <Tooltip>
              Volatility: 1%
              <span>
                The monthly metric moves up or down by 1% on average, after its
                first computation
              </span>
            </Tooltip>
          </BlackContainer>
        </SubSectionContainer>
        <SubSectionContainer>
          <strong>Current Drift</strong>
          <TansparentContainer>
            Detected on {new Date().toLocaleDateString()}
          </TansparentContainer>
          <TansparentContainer>7 month impacted</TansparentContainer>
          <BlackContainer>
            Total drift:<strong>48.9</strong>{" "}
            <span style={{ marginLeft: "auto", paddingLeft: "8px" }}>
              <TrendChip trend="up" absoluteValue={2} />
            </span>
          </BlackContainer>
        </SubSectionContainer>
        <SubSectionContainer>
          <strong>Owner</strong>
          <TansparentContainer>Aya Nakamura</TansparentContainer>
        </SubSectionContainer>
        <DrillDownButton
          onClick={() => {
            window.location.href =
              "/41231518/samox/dbt-example/overview?commitSha=37467fb6ce76d26fad8b09d7582ed3f6ad5d61e3&snapshotDate=2023-10-18";
          }}
        >
          <strong> DRILL DOWN</strong>
        </DrillDownButton>
      </DriftDetailContainer>
      <DualBarChartContainer>
        <DualMetricBarChart data={data} />
      </DualBarChartContainer>
    </PageContainer>
  );
};

export default DriftOverviewPage;
