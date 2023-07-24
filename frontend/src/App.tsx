import { useEffect, useState } from "react";
import "./App.css";
import { FunctionResponse, getCommitFiles } from "./services/github";

const [owner, repo, commitSHA] = [
  "Samox",
  "datadrift-example",
  "036f9d6b685ee02a14faa70ed05e0bd60650c477",
];

function App() {
  const [commitData, setCommitData] = useState<FunctionResponse | null>(null);

  useEffect(() => {
    const fetchCommitData = async () => {
      try {
        const response = await getCommitFiles(owner, repo, commitSHA);
        console.log("response", response);
        setCommitData(response);
      } catch (error) {
        console.error("Error fetching GitHub commit data:", error);
      }
    };

    fetchCommitData().catch(console.error);
  }, []);

  return (
    <>
      <a href={`https://github.com/${owner}/${repo}/commit/${commitSHA}`}>
        Link to commit {`${owner}/${repo}/commit/${commitSHA}`}
      </a>
      <div className="card">
        <p>The commit contains {commitData?.length || 0} file(s)</p>
      </div>
    </>
  );
}

export default App;
