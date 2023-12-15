import time
from datetime import datetime
from typing import Dict, List, Optional

import pandas as pd
from github import ContentFile, GithubException, Repository
from github.MainClass import Github

from ..dataframe.dataframe_update_breakdown import DataFrameUpdate
from ..drift_evaluator.drift_evaluators import drift_summary_to_string
from ..drift_evaluator.interface import DriftEvaluatorContext
from ..logger import get_logger
from .abstract_connector import AbstractConnector

logger = get_logger(__name__)


class GithubConnector(AbstractConnector):
    def __init__(
        self,
        github_client: Github,
        github_repository_name: str,
        default_branch: Optional[str] = None,
        assignees: Optional[List[str]] = None,
    ):
        self.repo = github_client.get_repo(github_repository_name)
        self.default_branch = default_branch if default_branch is not None else self.repo.default_branch
        self.assert_branch_exist(self.repo, self.default_branch)
        self.assignees = self.assert_assignees_exists(assignees if assignees is not None else [])
        self.logger = logger

    def assert_file_exists(self, file_path: str) -> Optional[ContentFile.ContentFile]:
        try:
            contents = self.repo.get_contents(file_path, ref=self.default_branch)
            assert not isinstance(contents, list), "pathfile returned multiple contents"
            return contents
        except GithubException as e:
            if e.status == 404:
                return None
            else:
                raise e

    @staticmethod
    def get_table_file_path(table_name: str) -> str:
        if table_name.endswith(".csv"):
            return table_name
        else:
            return table_name + ".csv"

    def get_table(self, table_name: str) -> Optional[pd.DataFrame]:
        file_path = self.get_table_file_path(table_name)
        table_file_content = self.assert_file_exists(file_path)
        if table_file_content is None:
            return None
        else:
            old_dataframe = pd.read_csv(
                table_file_content.download_url,
                dtype="string",
                keep_default_na=False,
            )
            return old_dataframe

    def init_table(self, table_name: str, dataframe: pd.DataFrame, measure_date: datetime):
        logger.info("Creating table on branch: " + self.default_branch)
        commit_message = "New data: " + table_name
        logger.info("Commit: " + commit_message)
        file_path = self.get_table_file_path(table_name)
        self.repo.create_file(
            file_path,
            commit_message,
            dataframe.to_csv(index=True, header=True),
            self.default_branch,
        )

    def open_issue(self, title: str, description_body: str):
        try:
            issue = self.repo.create_issue(
                title=title,
                body=description_body,
                assignees=self.assignees,
            )
            logger.info("Issue created: " + issue.html_url)
        except GithubException as e:
            if e.status == 422:
                logger.info("Issue already exists. skipping...")
            else:
                raise e

    @staticmethod
    def assert_branch_exist(repo: Repository.Repository, branch_name: str) -> None:
        try:
            branch = repo.get_branch(branch_name)
        except:
            branch = None

        if not branch:
            logger.info(f"Branch {branch_name} doesn't exist. Creating it...")
            reported_branch = repo.get_branch(repo.default_branch)

            repo.create_git_ref(f"refs/heads/{branch_name}", reported_branch.commit.sha)

    def assert_assignees_exists(self, maybe_assignees: List[str]) -> List[str]:
        members = [collaborator.login for collaborator in self.repo.get_collaborators()]
        exising_assignees = []
        for assignee in maybe_assignees:
            if assignee not in members:
                logger.warn(f"Assignee {assignee} does not exist")
            else:
                exising_assignees.append(assignee)
        return exising_assignees

    def checkout_branch_from_default_branch(self, branch_name: str):
        self.assert_branch_exist(self.repo, branch_name)
        try:
            ref = self.repo.get_git_ref(f"heads/{branch_name}")
            ref.delete()
        except GithubException:
            pass

        """
        Checkout a branch from the default branch of the given repository.
        """
        # Get the default branch of the repository
        default_branch = self.repo.get_branch(self.default_branch)
        logger.info(f"Checkout branch: {branch_name} from branch: {default_branch.name}")

        # Create a new reference to the default branch

        try:
            ref = self.repo.create_git_ref(f"refs/heads/{branch_name}", default_branch.commit.sha)
        except GithubException:
            pass
        return

    def update_file_with_retry(
        self,
        table_name: str,
        commit_message,
        data,
        branch,
        max_retries=3,
    ) -> str:
        retries = 0

        file_path = self.get_table_file_path(table_name)
        commit_sha = ""
        while retries < max_retries:
            try:
                content = self.assert_file_exists(file_path)
                if content is None:
                    response = self.repo.create_file(file_path, commit_message, data, branch)
                    logger.info(response["commit"].html_url)
                    commit_sha = response["commit"].sha
                else:
                    response = self.repo.update_file(file_path, commit_message, data, content.sha, branch)
                    logger.info(response["commit"].html_url)
                    commit_sha = response["commit"].sha
                return commit_sha
            except GithubException as e:
                if e.status == 409:
                    retries += 1
                    time.sleep(1)  # Wait for 1 second before retrying
                else:
                    raise e
        raise Exception(f"Failed to update ${table_name} file after {max_retries} retries")

    def handle_breakdown(
        self,
        table_name: str,
        update_breakdown: Dict[str, DataFrameUpdate],
        measure_date: datetime,
    ):
        branch = self.default_branch
        for key, value in update_breakdown.items():
            commit_message = f"{key}: {table_name}"
            if value.has_update:
                logger.info("Update: " + key)
                update_evaluation = value.update_evaluation
                update_context = value.update_context
                if update_context and update_evaluation:
                    commit_message += "\n\n" + update_evaluation.message
                    if isinstance(update_context, DriftEvaluatorContext) and update_context.summary != None:
                        summary = update_context.summary
                        drift_summary_string = drift_summary_to_string(summary)
                        commit_message += "\n\n" + drift_summary_string

                update_commit_sha = self.update_file_with_retry(
                    table_name=table_name,
                    commit_message=commit_message,
                    data=value.df.to_csv(index=True, header=True),
                    branch=branch,
                )
                if update_evaluation and update_evaluation.should_alert:
                    alert_message = "Commit: " + update_commit_sha
                    alert_message = alert_message + "\n\n" + update_evaluation.message
                    title = commit_message
                    self.open_issue(title=title, description_body=alert_message)
