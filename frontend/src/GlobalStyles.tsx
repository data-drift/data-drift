import { Global, css } from "@emotion/react";

export const GlobalStyles = () => (
  <Global
    styles={css`
      body {
        font-family: "JetBrains Mono", monospace;
      }
    `}
  />
);
