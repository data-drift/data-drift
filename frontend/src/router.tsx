import { createBrowserRouter } from "react-router-dom";
import GithubForm from "./pages/GithubForm";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <GithubForm />,
  },
]);
