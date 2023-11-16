import time
from typing import Optional, List
from datagit.dataframe_update_breakdown import (
    UpdateType,
    dataframe_update_breakdown,
)
import pandas as pd
from github import Github, Repository, ContentFile, GithubException
from datagit.drift_evaluators import (
    DefaultDriftEvaluator,
    DriftEvaluatorAbstractClass,
    drift_summary_to_string,
)
from datagit.dataset_helpers import (
    sort_dataframe_on_first_column_and_assert_is_unique,
)
import re
import os
import datetime


class GithubConnector:
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

    def init_file(
        self,
        file_path: str,
        dataframe: pd.DataFrame,
    ):
        commit_message = "New data: " + file_path
        print("Commit: " + commit_message)
        self.repo.create_file(
            file_path,
            commit_message,
            dataframe.to_csv(index=True, header=True),
            self.default_branch,
        )

    def create_pullrequest(self, file_path: str, description_body: str, branch: str):
        try:
            if len(self.assignees) > 0:
                pullrequest = self.repo.create_pull(
                    title="New drift detected " + file_path,
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


def store_table(
    *,
    github_client: Github,
    github_repository_name: str,
    branch: Optional[str] = None,
    assignees: Optional[List[str]] = None,
    table_dataframe: pd.DataFrame,
    table_name: str,
    drift_evaluator: DriftEvaluatorAbstractClass = DefaultDriftEvaluator(),
) -> None:
    """
    Store metrics into a specific repository file on GitHub.

    Parameters:
      ghClient (PyGithub.Github): An instance of a GitHub client to interact with the GitHub API.
      dataframe (pd.DataFrame): The dataframe containing the metrics to be stored.
      filepath (str): The full path to the target file in the format
        'organization/repository/path_to_file'.
      assignees (Optional[List[str]]): List of GitHub usernames to be assigned to the pull request.
        Defaults to None. If list is empty no alert will be raised, nor pull
        request will be created.
      branch (Optional[str]): The name of the branch where the metrics will be stored.
        If None, the default branch will be used. Defaults to None.
      drift_evaluator (Callable): Function that evaluates context and return information
        about how drift should be handled. See `drift_evaluator` module.

    Returns:
      None: This function does not return any value, but it performs a side effect of
      pushing the metric to GitHub.

    Raises:
      ValueError: If the dataframe does not have a unique first column or if the file
        path does not have the format 'organization/repository/path_to_file'.
      GithubException: If there is an error in interacting with the GitHub API, e.g.,
        insufficient permissions, non-existent repo, etc.
    """

    print("Storing metric...")
    drift_branch = get_alert_branch_name(table_name)

    github_connector = GithubConnector(
        github_client=github_client,
        github_repository_name=github_repository_name,
        default_branch=branch,
        assignees=assignees,
    )
    table_dataframe = sort_dataframe_on_first_column_and_assert_is_unique(
        table_dataframe
    )

    push_metric(
        table_dataframe,
        drift_branch,
        table_name,
        github_connector,
        drift_evaluator,
    )


def partition_and_store_table(
    *,
    github_client: Github,
    github_repository_name: str,
    table_dataframe: pd.DataFrame,
    table_name: str,
    branch: Optional[str] = None,
) -> None:
    """
    Store metrics into a specific repository file on GitHub.

    Parameters:
      ghClient (Github): An instance of a GitHub client to interact with the GitHub API.
      dataframe (pd.DataFrame): The dataframe containing the metrics to be stored.
      filepath (str): The full path to the target file in the format
        'organization/repository/path_to_file'.
      branch (Optional[str]): The name of the branch where the metrics will be stored.
        If None, the default branch will be used. Defaults to None.

    Returns:
      None: This function does not return any value, but it performs a side effect of
      pushing the metric to GitHub.

    Raises:
      ValueError: If the dataframe does not have a unique first column or if the file
        path does not have the format 'organization/repository/path_to_file'.
      GithubException: If there is an error in interacting with the GitHub API, e.g.,
        insufficient permissions, non-existent repo, etc.
    """

    print("Partitionning metric...")

    table_dataframe["date"] = pd.to_datetime(table_dataframe["date"])

    grouped = table_dataframe.groupby(pd.Grouper(key="date", freq="M"))

    # Iterate over the groups and print the sub-dataframes
    for name, group in grouped:
        print(f"Storing metric for Month: {name}")
        monthly_table_name = get_monthly_file_path(table_name, name.strftime("%Y-%m"))  # type: ignore
        store_table(
            github_client=github_client,
            table_dataframe=group,
            github_repository_name=github_repository_name,
            table_name=monthly_table_name,
            branch=branch,
            assignees=None,
        )


def push_metric(
    dataframe,
    drift_branch,
    file_path,
    github_connector: GithubConnector,
    drift_evaluator: DriftEvaluatorAbstractClass = DefaultDriftEvaluator(),
):
    default_branch = github_connector.default_branch
    if dataframe.index.name != "unique_key":
        dataframe = dataframe.set_index("unique_key")

    dataframe = dataframe.astype("string")
    contents = github_connector.assert_file_exists(file_path)
    if contents is None:
        print("Metric not found, creating it on branch: " + default_branch)
        github_connector.init_file(file_path=file_path, dataframe=dataframe)
        print("Metric stored")
        pass
    else:
        print("Metric found, updating it on branch: " + default_branch)
        date_column = find_date_column(dataframe)
        if contents.content is not None and date_column is not None:
            # Compare the contents of the file with the new contents and assert if it need 2 commits
            print("Content", contents.download_url)
            print("Dataframe dtypes", dataframe.dtypes.to_dict())
            old_dataframe = pd.read_csv(
                contents.download_url,
                dtype="string",
                keep_default_na=False,
            )
            print("Old Dataframe dtypes", old_dataframe.dtypes.to_dict())
            update_breakdown = dataframe_update_breakdown(
                old_dataframe, dataframe, drift_evaluator
            )
            if any(item["has_update"] for item in update_breakdown.values()):
                print("Change detected")
            else:
                print("Nothing to update")
                pass
            branch = default_branch
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
                        drift_summary_string = ""
                        if value["drift_summary"]:
                            drift_summary_string = drift_summary_to_string(
                                value["drift_summary"]
                            )
                            commit_message += "\n\n" + drift_summary_string
                        if drift_evaluation["should_alert"]:
                            if branch == default_branch:
                                github_connector.checkout_branch_from_default_branch(
                                    drift_branch
                                )
                                branch = drift_branch
                            pr_message = (
                                pr_message + "\n\n" + drift_evaluation["message"]
                            )

                    github_connector.update_file_with_retry(
                        file_path=file_path,
                        commit_message=commit_message,
                        data=value["df"].to_csv(index=True, header=True),
                        branch=branch,
                    )

            if pr_message != "":
                github_connector.create_pullrequest(file_path, pr_message, branch)


def get_alert_branch_name(filepath: str) -> str:
    """
    Returns a valid Git branch name based on the given filepath.
    """
    # Replace any non-alphanumeric characters with hyphens
    branch_name = re.sub(r"[^a-zA-Z0-9]+", "-", filepath)

    # Remove any leading or trailing hyphens
    branch_name = branch_name.strip("-")

    # Convert to lowercase
    branch_name = branch_name.lower()
    now = datetime.datetime.now()
    datetime_str = now.strftime("%Y-%m-%d-%H-%M-%S")
    # Append a prefix
    branch_name = f"drift/{datetime_str}/{branch_name}"

    # Truncate to 63 characters (the maximum allowed length for a Git branch name)
    branch_name = branch_name[:63]

    return branch_name


def find_date_column(df):
    date_columns = df.filter(like="date").columns
    if len(date_columns) > 0:
        return date_columns[0]
    else:
        return df.columns[0]


def get_monthly_file_path(file_path, month):
    directory, file_name = os.path.split(file_path)
    file_name, extension = os.path.splitext(file_name)

    new_file_name = f"{file_name}/{month}{extension}"

    new_file_path = os.path.join(directory, new_file_name)

    return new_file_path
