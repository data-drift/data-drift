import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import { GlobalStyles } from "../../GlobalStyles.tsx";
import { ThemeProvider } from "@emotion/react";
import { theme } from "../../theme.ts";
import { SnapshotDiff } from "./generatedDiff.ts";

// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ThemeProvider theme={theme}>
      <GlobalStyles />
      <App />
    </ThemeProvider>
  </React.StrictMode>
);

declare global {
  interface Window {
    generated_diff: SnapshotDiff;
  }
}
