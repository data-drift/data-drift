import { keyframes } from "@emotion/react";
import styled from "@emotion/styled";

const spin = keyframes`
  0% { transform: rotate(0deg); }
  25% { transform: rotate(110deg); }
  50% { transform: rotate(190deg); }
  75% { transform: rotate(295deg); }
  100% { transform: rotate(360deg); }
`;

const LoaderWrapper = styled.div`
  background-color: #121212;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  overflow: hidden;
`;

const BrutalistLoader = styled.div`
  position: relative;
  width: 100px;
  height: 100px;
`;

const Circle = styled.div`
  width: 100%;
  height: 100%;
  border: 5px solid grey;
  border-top-color: white;
  border-radius: 50%;
  animation: ${spin} 1s linear infinite;
`;

const Loader = () => {
  return (
    <LoaderWrapper>
      <BrutalistLoader>
        <Circle />
      </BrutalistLoader>
    </LoaderWrapper>
  );
};

export default Loader;
