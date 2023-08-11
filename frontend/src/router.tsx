import { createBrowserRouter } from "react-router-dom";
import GithubForm from "./pages/GithubForm";
import DisplayCommit from "./pages/DisplayCommit";
import MetricCohort from "./pages/MetricCohorts";

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
  {
    path: "report/:installationId/metrics/:metricName/cohorts/:timegrain",
    element: <MetricCohort />,
    loader: MetricCohort.loader,
  },
]);
