import { FormEvent, useState, useEffect } from "react";
import { LOCAL_STORAGE_GITHUB_TOKEN, parseGithubUrl } from "./services/github";

function GithubForm() {
  const [url, setUrl] = useState("");
  const [token, setToken] = useState("");

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    localStorage.setItem(LOCAL_STORAGE_GITHUB_TOKEN, token);
    const commitInfo = parseGithubUrl(url);
    console.log(commitInfo);
    window.location.href = `/${commitInfo.owner}/${commitInfo.repo}/commit/${commitInfo.commitSHA}`;
  };

  useEffect(() => {
    const storedToken = localStorage.getItem(LOCAL_STORAGE_GITHUB_TOKEN);
    if (storedToken) {
      setToken(storedToken);
    }
  }, []);

  return (
    <form onSubmit={handleSubmit}>
      <label>
        <span> GitHub Url:</span>
        <input
          type="text"
          value={url}
          style={{ width: "300px" }}
          onChange={(e) => setUrl(e.target.value)}
        />
      </label>
      <br />
      <label>
        <span>GitHub Token:</span>

        <input
          type="text"
          value={token}
          style={{ width: "300px" }}
          onChange={(e) => setToken(e.target.value)}
        />
      </label>
      <br />
      <button type="submit">Submit</button>
    </form>
  );
}

export default GithubForm;
