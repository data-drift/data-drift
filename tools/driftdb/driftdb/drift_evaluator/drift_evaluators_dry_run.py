import traceback

import pandas as pd
from driftdb.dataframe.summarize_dataframe_updates import summarize_dataframe_updates
from github.MainClass import Github

from ..logger import get_logger
from .drift_evaluators import BaseDriftEvaluator, BaseNewDataEvaluator, DriftEvaluatorContext, parse_drift_summary
from .interface import DriftEvaluation, NewDataEvaluatorContext

logger = get_logger(__name__)


def run_drift_evaluator(
    *, drift_evaluator: BaseDriftEvaluator, gh_client: Github, repo_name: str, commit_sha: str
) -> DriftEvaluation:
    repo = gh_client.get_repo(repo_name)
    commit = repo.get_commit(commit_sha)
    file = commit.files[0]
    raw_content = repo.get_contents(file.filename, ref=commit_sha)
    parent_raw_content = repo.get_contents(file.filename, ref=commit.parents[0].sha)

    if isinstance(raw_content, list):
        raw_content = raw_content[0]

    if isinstance(parent_raw_content, list):
        parent_raw_content = parent_raw_content[0]

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

    drift_summary = summarize_dataframe_updates(initial_df=old_dataframe, final_df=new_dataframe)

    #  run drift evaluator
    data_drift_context = DriftEvaluatorContext(
        before=old_dataframe,
        after=new_dataframe,
        summary=drift_summary,
    )
    drift_evaluation = drift_evaluator.compute_drift_evaluation(data_drift_context)
    return drift_evaluation


def run_new_data_evaluator(
    *, drift_evaluator: BaseNewDataEvaluator, gh_client: Github, repo_name: str, commit_sha: str
) -> DriftEvaluation:
    repo = gh_client.get_repo(repo_name)
    commit = repo.get_commit(commit_sha)
    file = commit.files[0]
    raw_content = repo.get_contents(file.filename, ref=commit_sha)
    parent_raw_content = repo.get_contents(file.filename, ref=commit.parents[0].sha)

    if isinstance(raw_content, list):
        raw_content = raw_content[0]

    if isinstance(parent_raw_content, list):
        parent_raw_content = parent_raw_content[0]

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

    new_data = new_dataframe.loc[~new_dataframe["unique_key"].isin(old_dataframe["unique_key"])]

    #  run drift evaluator
    data_drift_context = NewDataEvaluatorContext(
        before=old_dataframe,
        after=new_dataframe,
        added_rows=new_data,
    )
    drift_evaluation = drift_evaluator.compute_new_data_evaluation(data_drift_context)
    return drift_evaluation
