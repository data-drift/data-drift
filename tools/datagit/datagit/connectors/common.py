import datetime
import os
import re
from typing import Optional, List
from .github_connector import GithubConnector
from ..dataframe.dataframe_update_breakdown import (
    UpdateType,
    dataframe_update_breakdown,
)
from ..dataframe.helpers import (
    sort_dataframe_on_first_column_and_assert_is_unique,
)
from ..drift_evaluator.drift_evaluators import (
    DefaultDriftEvaluator,
    DriftEvaluatorAbstractClass,
    drift_summary_to_string,
)
from github import Github

import pandas as pd


def snapshot_table(
    table_dataframe: pd.DataFrame,
    table_name: str,
    github_connector: GithubConnector,
    drift_evaluator: DriftEvaluatorAbstractClass = DefaultDriftEvaluator(),
):
    table_dataframe = sort_dataframe_on_first_column_and_assert_is_unique(
        table_dataframe
    )
    if table_dataframe.index.name != "unique_key":
        table_dataframe = table_dataframe.set_index("unique_key")
    table_dataframe = table_dataframe.astype("string")

    default_branch = github_connector.default_branch

    latest_stored_snapshot = github_connector.get_latest_table_snapshot(table_name)

    if latest_stored_snapshot is None:
        print("Table not found, creating it on branch: " + default_branch)
        github_connector.init_file(file_path=table_name, dataframe=table_dataframe)
        print("Table stored")
        pass
    else:
        print("Table found, updating it on branch: " + default_branch)
        date_column = find_date_column(table_dataframe)
        if date_column is not None:
            # Compare the contents of the file with the new contents and assert if it need 2 commits
            print("Dataframe dtypes", table_dataframe.dtypes.to_dict())

            print("Old Dataframe dtypes", latest_stored_snapshot.dtypes.to_dict())
            update_breakdown = dataframe_update_breakdown(
                latest_stored_snapshot, table_dataframe, drift_evaluator
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
                                drift_branch = get_alert_branch_name(table_name)
                                github_connector.checkout_branch_from_default_branch(
                                    drift_branch
                                )
                                branch = drift_branch
                            pr_message = (
                                pr_message + "\n\n" + drift_evaluation["message"]
                            )

                    github_connector.update_file_with_retry(
                        file_path=table_name,
                        commit_message=commit_message,
                        data=value["df"].to_csv(index=True, header=True),
                        branch=branch,
                    )

            if pr_message != "":
                github_connector.create_pullrequest(table_name, pr_message, branch)


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
    Store tables into a specific repository file on GitHub.

    Parameters:
      ghClient (PyGithub.Github): An instance of a GitHub client to interact with the GitHub API.
      dataframe (pd.DataFrame): The dataframe containing the tables to be stored.
      filepath (str): The full path to the target file in the format
        'organization/repository/path_to_file'.
      assignees (Optional[List[str]]): List of GitHub usernames to be assigned to the pull request.
        Defaults to None. If list is empty no alert will be raised, nor pull
        request will be created.
      branch (Optional[str]): The name of the branch where the tables will be stored.
        If None, the default branch will be used. Defaults to None.
      drift_evaluator (Callable): Function that evaluates context and return information
        about how drift should be handled. See `drift_evaluator` module.

    Returns:
      None: This function does not return any value, but it performs a side effect of
      pushing the table to GitHub.

    Raises:
      ValueError: If the dataframe does not have a unique first column or if the file
        path does not have the format 'organization/repository/path_to_file'.
      GithubException: If there is an error in interacting with the GitHub API, e.g.,
        insufficient permissions, non-existent repo, etc.
    """

    print("Storing table...")

    github_connector = GithubConnector(
        github_client=github_client,
        github_repository_name=github_repository_name,
        default_branch=branch,
        assignees=assignees,
    )

    snapshot_table(
        table_dataframe,
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
    Store tables into a specific repository file on GitHub.

    Parameters:
      ghClient (Github): An instance of a GitHub client to interact with the GitHub API.
      dataframe (pd.DataFrame): The dataframe containing the tables to be stored.
      filepath (str): The full path to the target file in the format
        'organization/repository/path_to_file'.
      branch (Optional[str]): The name of the branch where the tables will be stored.
        If None, the default branch will be used. Defaults to None.

    Returns:
      None: This function does not return any value, but it performs a side effect of
      pushing the table to GitHub.

    Raises:
      ValueError: If the dataframe does not have a unique first column or if the file
        path does not have the format 'organization/repository/path_to_file'.
      GithubException: If there is an error in interacting with the GitHub API, e.g.,
        insufficient permissions, non-existent repo, etc.
    """

    github_connector = GithubConnector(
        github_client=github_client,
        github_repository_name=github_repository_name,
        default_branch=branch,
    )

    print("Partitionning table by month...")

    table_dataframe["date"] = pd.to_datetime(table_dataframe["date"])

    grouped = table_dataframe.groupby(pd.Grouper(key="date", freq="M"))

    # Iterate over the groups and print the sub-dataframes
    for name, group in grouped:
        print(f"Storing table for Month: {name}")
        monthly_table_name = get_monthly_file_path(table_name, name.strftime("%Y-%m"))  # type: ignore
        snapshot_table(
            table_dataframe=group,
            table_name=monthly_table_name,
            github_connector=github_connector,
        )


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
