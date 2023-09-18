import { ddCommitDiffUrlFactory } from "../../services/data-drift";
import { CommitListItem } from "./CommitListItem";

type CommitDataFromApi = {
  sha: string;
  commit: {
    author: {
      date?: string | undefined;
    } | null;
    message: string;
  };
}[];

export const CommitList = ({
  data,
  params,
}: {
  data: CommitDataFromApi;
  params: {
    installationId: string;
    owner: string;
    repo: string;
  };
}) => {
  console.log("data", data);
  return (
    <div>
      {data.map((commit) => {
        const isDrift = commit.commit.message.includes("Drift");
        const commitUrl = ddCommitDiffUrlFactory({
          ...params,
          commitSha: commit.sha,
        });
        return (
          <CommitListItem
            key={commit.sha}
            type={isDrift ? "Drift" : "New Data"}
            date={
              commit.commit.author?.date
                ? new Date(commit.commit.author?.date)
                : null
            }
            name={commit.commit.message}
            commitUrl={commitUrl}
          />
        );
      })}
    </div>
  );
};
