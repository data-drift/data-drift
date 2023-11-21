import traceback

import pandas as pd
from github import Github

from ..logger import get_logger
from .drift_evaluators import BaseDriftEvaluator, DriftEvaluatorContext, parse_drift_summary

logger = get_logger(__name__)


def run_drift_evaluator(*, drift_evaluator: BaseDriftEvaluator, gh_client: Github, repo_name: str, commit_sha: str):
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

    drift_summary = None
    try:
        commit_message = commit.commit.message
        drift_summary = parse_drift_summary(commit_message)
    except Exception as e:
        logger.warn("Failed to parse drift summary: " + str(e))
        traceback.print_exc()

    #  run drift evaluator
    data_drift_context = DriftEvaluatorContext(
        before=old_dataframe,
        after=new_dataframe,
        summary=drift_summary,
    )
    try:
        drift_evaluation = drift_evaluator.compute_drift_evaluation(data_drift_context)
        #  return result
        return drift_evaluation
    except Exception as e:
        logger.warn("Drift evaluator failed: " + str(e))
        traceback.print_exc()
        logger.warn("Drift evaluator failed: " + str(e))
        traceback.print_exc()
