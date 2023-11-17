from datetime import datetime
import time
from typing import Dict, Optional, List

from .abstract_connector import AbstractConnector
from .common import get_alert_branch_name

from ..drift_evaluator.drift_evaluators import drift_summary_to_string
from ..dataframe.dataframe_update_breakdown import DataFrameUpdate, UpdateType
import pandas as pd
from github import Github, Repository, ContentFile, GithubException


class GithubConnector(AbstractConnector):
    def __init__(
        self,
        github_client: Github,
        github_repository_name: str,
        default_branch: Optional[str] = None,
        assignees: Optional[List[str]] = None,
    ):
        self.repo = github_client.get_repo(github_repository_name)
        self.default_branch = (
            default_branch if default_branch is not None else self.repo.default_branch
        )
        self.assert_branch_exist(self.repo, self.default_branch)
        self.assignees = assignees if assignees is not None else []

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

    def get_table(self, table_name: str) -> Optional[pd.DataFrame]:
        table_file_content = self.assert_file_exists(table_name)
        if table_file_content is None:
            return None
        else:
            old_dataframe = pd.read_csv(
                table_file_content.download_url,
                dtype="string",
                keep_default_na=False,
            )
            return old_dataframe

    def init_table(
        self, table_name: str, dataframe: pd.DataFrame, measure_date: datetime
    ):
        print("Creating table on branch: " + self.default_branch)
        commit_message = "New data: " + table_name
        print("Commit: " + commit_message)
        self.repo.create_file(
            table_name,
            commit_message,
            dataframe.to_csv(index=True, header=True),
            self.default_branch,
        )

    def close_pullrequests(self, title: str):
        pulls = self.repo.get_pulls(state="open")
        for pull in pulls:
            if pull.title == title:
                pull.edit(state="closed")

    def create_pullrequest(self, title: str, description_body: str, branch: str):
        try:
            if len(self.assignees) > 0:
                pullrequest = self.repo.create_pull(
                    title=title,
                    body=description_body,
                    head=branch,
                    base=self.default_branch,
                )
                print("Pull request created: " + pullrequest.html_url)
                existing_assignees = self.assert_assignees_exists()
                pullrequest.add_to_assignees(*existing_assignees)
            else:
                print("No assignees, skipping pull request creation")
        except GithubException as e:
            if e.status == 422:
                print("Pull request already exists, skipping...")
            else:
                raise e

    @staticmethod
    def assert_branch_exist(repo: Repository.Repository, branch_name: str) -> None:
        try:
            branch = repo.get_branch(branch_name)
        except:
            branch = None

        if not branch:
            print(f"Branch {branch_name} doesn't exist, creating it...")
            reported_branch = repo.get_branch(repo.default_branch)

            repo.create_git_ref(f"refs/heads/{branch_name}", reported_branch.commit.sha)

    def assert_assignees_exists(self) -> List[str]:
        members = [collaborator.login for collaborator in self.repo.get_collaborators()]
        exising_assignees = []
        for assignee in self.assignees:
            if assignee not in members:
                print(f"Assignee {assignee} does not exist")
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
        print(
            "Checkout branch: " + branch_name,
            " from branch:" + default_branch.name,
        )

        # Create a new reference to the default branch

        try:
            ref = self.repo.create_git_ref(
                f"refs/heads/{branch_name}", default_branch.commit.sha
            )
        except GithubException:
            pass
        return

    def update_file_with_retry(
        self,
        file_path,
        commit_message,
        data,
        branch,
        max_retries=3,
    ):
        retries = 0

        while retries < max_retries:
            try:
                content = self.assert_file_exists(file_path)
                if content is None:
                    response = self.repo.create_file(
                        file_path, commit_message, data, branch
                    )
                    print(response["commit"].html_url)
                else:
                    response = self.repo.update_file(
                        file_path, commit_message, data, content.sha, branch
                    )
                    print(response["commit"].html_url)
                return
            except GithubException as e:
                if e.status == 409:
                    retries += 1
                    time.sleep(1)  # Wait for 1 second before retrying
                else:
                    raise e
        raise Exception(f"Failed to update file after {max_retries} retries")

    def handle_breakdown(
        self,
        table_name: str,
        update_breakdown: Dict[str, DataFrameUpdate],
        measure_date: datetime,
    ):
        branch = self.default_branch
        pr_message = ""
        for key, value in update_breakdown.items():
            commit_message = key
            if value["has_update"]:
                print("Update: " + key)
                if (
                    value["type"] == UpdateType.DRIFT
                    and value["drift_context"]
                    and value["drift_evaluation"]
                ):
                    drift_evaluation = value["drift_evaluation"]
                    commit_message += "\n\n" + drift_evaluation["message"]
                    if value["drift_summary"]:
                        drift_summary_string = drift_summary_to_string(
                            value["drift_summary"]
                        )
                        commit_message += "\n\n" + drift_summary_string
                    if drift_evaluation["should_alert"]:
                        if branch == self.default_branch:
                            drift_branch = get_alert_branch_name(table_name)
                            self.checkout_branch_from_default_branch(drift_branch)
                            branch = drift_branch
                        pr_message = pr_message + "\n\n" + drift_evaluation["message"]

                self.update_file_with_retry(
                    file_path=table_name,
                    commit_message=commit_message,
                    data=value["df"].to_csv(index=True, header=True),
                    branch=branch,
                )

        if pr_message != "":
            title = "New drift detected " + table_name
            self.close_pullrequests(title=title)
            self.create_pullrequest(
                title=title,
                description_body=pr_message,
                branch=branch,
            )
