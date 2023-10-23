import styled from "@emotion/styled";

export const Tooltip = styled.span`
  position: relative;
  cursor: pointer;

  &:hover > span {
    display: block;
  }

  & > span {
    display: none;
    position: absolute;
    top: 50%;
    left: 100%; // Position it to the right of the parent
    transform: translateY(-50%); // Center it vertically
    padding: 8px;
    background-color: ${({ theme }) => theme.colors.background};
    border-radius: 0;
    border: "1px solid";
    border-color: ${({ theme }) => theme.colors.background2};
    color: white;
    font-size: 0.8em;
    white-space: normal;
    width: 200px;
    z-index: 1000;
  }
`;
