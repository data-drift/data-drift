import React from "react";
import ReactDOM from "react-dom";
import App from "./App.tsx";
import { GlobalStyles } from "../../GlobalStyles.tsx";
import { ThemeProvider } from "@emotion/react";
import { theme } from "../../theme.ts";

// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ThemeProvider theme={theme}>
      <GlobalStyles />
      <App />
    </ThemeProvider>
  </React.StrictMode>
);
