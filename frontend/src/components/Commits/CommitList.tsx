import { useMemo } from "react";
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
  filters,
}: {
  data: CommitDataFromApi;
  params: {
    installationId: string;
    owner: string;
    repo: string;
  };
  filters: {
    driftDate: {
      driftDate: string;
      isChecked: boolean;
    };
    parentData: {
      parentData: string[];
      childChecked: boolean[];
    };
  };
}) => {
  const filteredData = useMemo(() => {
    const filteredData = filters.driftDate.isChecked
      ? data.filter((commit) => {
          const date = new Date(commit.commit.author?.date || "");
          const driftDate = new Date(filters.driftDate.driftDate);
          const diff = Math.abs(date.getTime() - driftDate.getTime());
          const hoursDiff = diff / (1000 * 60 * 60);
          return hoursDiff <= 12;
        })
      : data;
    return filteredData;
  }, [data, filters]);
  return (
    <div>
      {filteredData.length > 0 ? (
        filteredData.map((commit) => {
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
        })
      ) : (
        <div
          style={{
            border: "1px solid #ccc",
            padding: "16px",
            borderRadius: "0",
            marginBottom: "16px",
            display: "flex",
            flexDirection: "column",
            alignItems: "flex-start",
            width: "100%",
          }}
        >
          No commits found
        </div>
      )}
    </div>
  );
};
