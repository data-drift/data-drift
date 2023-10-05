import { createBrowserRouter } from "react-router-dom";
import GithubForm from "./pages/GithubForm";
import DisplayCommit from "./pages/DisplayCommit/DisplayCommit";
import MetricCohort from "./pages/MetricCohorts";
import MetricReportWaterfall from "./pages/MetricReportWaterfall";
import { HomePage } from "./pages/HomePage";
import DriftListPage from "./pages/DriftList";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <HomePage />,
  },
  {
    path: "/ghform",
    element: <GithubForm />,
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
  {
    path: "report/:installationId/metrics/:metricName/report/:timegrainValue",
    element: <MetricReportWaterfall />,
    loader: MetricReportWaterfall.loader,
  },
  {
    path: "report/:installationId/:owner/:repo/commits",
    element: <DriftListPage />,
    loader: DriftListPage.loader,
  },
]);
