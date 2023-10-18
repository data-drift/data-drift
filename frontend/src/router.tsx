import { createBrowserRouter } from "react-router-dom";
import GithubForm from "./pages/GithubForm";
import DisplayCommit from "./pages/DisplayCommit/DisplayCommit";
import MetricCohort from "./pages/MetricCohorts";
import MetricReportWaterfall from "./pages/MetricReportWaterfall";
import { HomePage } from "./pages/HomePage";
import DriftListPage from "./pages/DriftList";
import Overview from "./pages/Overview/Overview";
import TableList from "./pages/TableList";
import TablePage from "./pages/TablePage";
import MetricPage from "./pages/MetricPage";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <HomePage />,
  },
  {
    path: "/:installationId/:owner/:repo/overview",
    element: <Overview />,
    loader: Overview.loader,
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
  {
    path: "tables",
    element: <TableList />,
    loader: TableList.loader,
  },
  {
    path: "tables/:tableName",
    element: <TablePage />,
    loader: TablePage.loader,
  },
  {
    path: "tables/:tableName/metrics/:metricName",
    element: <MetricPage />,
    loader: MetricPage.loader,
  },
]);
