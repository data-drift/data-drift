import { FormEvent, useState, useEffect } from "react";
import { LOCAL_STORAGE_GITHUB_TOKEN } from "./services/github";

function GithubForm() {
  const [url, setUrl] = useState("");
  const [token, setToken] = useState("");

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    localStorage.setItem(LOCAL_STORAGE_GITHUB_TOKEN, token);
    window.location.href = url;
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
        URL:
        <input
          type="text"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
        />
      </label>
      <br />
      <label>
        GitHub Token:
        <input
          type="text"
          value={token}
          onChange={(e) => setToken(e.target.value)}
        />
      </label>
      <br />
      <button type="submit">Submit</button>
    </form>
  );
}

export default GithubForm;
