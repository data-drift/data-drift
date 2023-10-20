import styled from "@emotion/styled";

const PageContainer = styled.div`
  height: 100vh;
  width: 100vw;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  padding-left: ${({ theme }) => theme.spacing(6)};
  align-items: flex-start;
`;

const Separator = styled.div`
  width: 100%;
  border-bottom: 1px solid ${({ theme }) => theme.colors.text};
  margin-bottom: ${({ theme }) => theme.spacing(2)};
`;

const DriftOverviewPage = () => {
  return (
    <PageContainer>
      <h1>Drift overview</h1>
      <Separator />
    </PageContainer>
  );
};

export default DriftOverviewPage;
