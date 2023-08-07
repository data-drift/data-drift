import { createBrowserRouter } from "react-router-dom";
import GithubForm from "./pages/GithubForm";
import DisplayCommit from "./pages/DisplayCommit";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <GithubForm />,
  },
  {
    path: "/:owner/:repo/commit/:commitSHA",
    element: <DisplayCommit />,
    loader: DisplayCommit.githubLoader,
  },
  {
    path: "report/:installationId/:owner/:repo/commit/:commitSHA",
    element: <DisplayCommit />,
    loader: DisplayCommit.dataDriftLoader,
  },
]);
