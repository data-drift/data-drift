import { Global, css, useTheme } from "@emotion/react";

export const GlobalStyles = () => {
  const theme = useTheme();
  return (
    <Global
      styles={css`
        body {
          font-family: "JetBrains Mono", monospace;
          background-color: ${theme.colors.background2};
          color: ${theme.colors.text};
        }
      `}
    />
  );
};
