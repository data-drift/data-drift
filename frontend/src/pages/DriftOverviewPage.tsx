import styled from "@emotion/styled";

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
`;

const DriftOverviewPage = () => {
  return (
    <PageContainer>
      <h1>Drift overview</h1>
      <Separator />
      <DriftDetailContainer>Coucou</DriftDetailContainer>
    </PageContainer>
  );
};

export default DriftOverviewPage;
