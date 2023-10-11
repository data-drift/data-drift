import { LineageEvent } from "../../components/Lineage/MetricNode";
import { getCommitList } from "../../services/data-drift";
import { getFileCommits } from "./flow-nodes";

describe("getFileCommits", () => {
  const commitList = [
    {
      sha: "commit1",
      commit: { message: "Drift: path/to/file.csv" },
    },
    {
      sha: "commit2",
      commit: { message: "New data: path/to/file.csv" },
    },
    {
      sha: "commit3",
      commit: { message: "New data: path/to/other/file/2022-08.csv" },
    },
    {
      sha: "commit4",
      commit: { message: "New data: path/to/other/file/2022-09.csv" },
    },
    {
      sha: "commit5",
      commit: { message: "New data: path/to/other/file/2022-10.csv" },
    },
    {
      sha: "commit6",
      commit: { message: "New data: another/file/2022-08.csv" },
    },
    {
      sha: "commit7",
      commit: { message: "New data: another/file/2022-09.csv" },
    },
    {
      sha: "commit8",
      commit: { message: "Drift: another/file/2022-10.csv" },
    },
  ] as unknown as Awaited<ReturnType<typeof getCommitList>>["data"];

  it("should return an array of lineage events for a file", () => {
    const filepath = "path/to/file";
    const selectCommit = jest.fn<(commitSha: string) => void, [string]>();
    const expectedEvents = [
      { type: "Drift", onClick: expect.any(Function) as () => void },
      { type: "New Data", onClick: expect.any(Function) as () => void },
    ] satisfies LineageEvent[];

    const events = getFileCommits(commitList, filepath, selectCommit);

    expect(events).toEqual(expectedEvents);
  });

  it("should return an empty array if the file path is not found", () => {
    const filepath = "non-existent-file";
    const selectCommit = jest.fn();

    const events = getFileCommits(commitList, filepath, selectCommit);

    expect(events).toEqual([]);
    expect(selectCommit).not.toHaveBeenCalled();
  });

  it("should group commit of the same type when partition", () => {
    const filepath = "path/to/other/file";
    const selectCommit = jest.fn();
    const expectedEvents = [
      {
        type: "New Data",
        onClick: undefined,
        subEvents: [
          { name: "2022-08", onClick: expect.any(Function) as () => void },
          { name: "2022-09", onClick: expect.any(Function) as () => void },
          { name: "2022-10", onClick: expect.any(Function) as () => void },
        ],
      },
    ];

    const events = getFileCommits(commitList, filepath, selectCommit);

    expect(events).toEqual(expectedEvents);
  });
  it("should group commit of the same type when partition another file", () => {
    const filepath = "another/file";
    const selectCommit = jest.fn();
    const expectedEvents = [
      {
        type: "New Data",
        onClick: undefined,
        subEvents: [
          { name: "2022-08", onClick: expect.any(Function) as () => void },
          { name: "2022-09", onClick: expect.any(Function) as () => void },
        ],
      },
      {
        type: "Drift",
        onClick: undefined,
        subEvents: [
          { name: "2022-10", onClick: expect.any(Function) as () => void },
        ],
      },
    ];

    const events = getFileCommits(commitList, filepath, selectCommit);

    expect(events).toEqual(expectedEvents);
  });
});
