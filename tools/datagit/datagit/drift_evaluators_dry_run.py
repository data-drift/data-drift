from typing import Callable, Dict
from github import Github

import pandas as pd


def run_drift_evaluator(
    *,
    drift_evaluator: Callable[[Dict[str, pd.DataFrame]], Dict],
    gh_client: Github,
    repo_name: str,
    commit_sha: str
):
    #  get drift context from gh_client and commit_sha
    repo = gh_client.get_repo(repo_name)
    commit = repo.get_commit(commit_sha)
    file = commit.files[0]
    raw_content = repo.get_contents(file.filename, ref=commit_sha)
    parent_raw_content = repo.get_contents(file.filename, ref=commit.parents[0].sha)

    new_dataframe = pd.read_csv(
        raw_content.download_url,
        dtype="string",
        keep_default_na=False,
    )

    old_dataframe = pd.read_csv(
        parent_raw_content.download_url,
        dtype="string",
        keep_default_na=False,
    )

    #  run drift evaluator
    data_drift_context = {
        "reported_dataframe": old_dataframe,
        "computed_dataframe": new_dataframe,
    }
    drift_evaluation = drift_evaluator(data_drift_context)
    #  return result
    return drift_evaluation
