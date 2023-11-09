import time
from typing import Optional, List, Callable, Dict
from datagit.dataframe_update_breakdown import UpdateType, dataframe_update_breakdown
import pandas as pd
from github import Github, Repository, ContentFile, GithubException
from datagit.drift_evaluators import (
    DriftEvaluatorContext,
    auto_merge_drift,
    safe_drift_evaluator,
)
from datagit.dataset_helpers import (
    sort_dataframe_on_first_column_and_assert_is_unique,
)
import re
import os
import datetime


def store_metric(
    *,
    ghClient: Github,
    dataframe: pd.DataFrame,
    filepath: str,
    branch: Optional[str] = None,
    assignees: Optional[List[str]] = None,
    drift_evaluator: Callable[[DriftEvaluatorContext], Dict] = auto_merge_drift,
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
    if assignees is None:
        assignees = []

    print("Storing metric...")
    repo_orga, repo_name, file_path = filepath.split("/", 2)
    drift_branch = get_valid_branch_name(file_path)

    repo = ghClient.get_repo(repo_orga + "/" + repo_name)
    working_branch = branch if branch is not None else repo.default_branch
    assert_branch_exist(repo, working_branch)
    dataframe = sort_dataframe_on_first_column_and_assert_is_unique(dataframe)

    push_metric(
        dataframe,
        assignees,
        working_branch,
        drift_branch,
        file_path,
        repo,
        drift_evaluator,
    )


def partition_and_store_table(
    *,
    ghClient: Github,
    dataframe: pd.DataFrame,
    filepath: str,
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

    dataframe["date"] = pd.to_datetime(dataframe["date"])

    grouped = dataframe.groupby(pd.Grouper(key="date", freq="M"))

    # Iterate over the groups and print the sub-dataframes
    for name, group in grouped:
        print(f"Storing metric for Month: {name}")
        monthly_filepath = get_monthly_file_path(filepath, name.strftime("%Y-%m"))  # type: ignore
        store_metric(
            ghClient=ghClient,
            dataframe=group,
            filepath=monthly_filepath,
            branch=branch,
            assignees=None,
            drift_evaluator=auto_merge_drift,
        )


def push_metric(
    dataframe,
    assignees,
    default_branch,
    drift_branch,
    file_path,
    repo,
    drift_evaluator: Callable[[DriftEvaluatorContext], Dict],
):
    if dataframe.index.name != "unique_key":
        dataframe = dataframe.set_index("unique_key")

    dataframe = dataframe.astype("string")
    contents = assert_file_exists(repo, file_path, ref=default_branch)
    if contents is None:
        print("Metric not found, creating it on branch: " + default_branch)
        create_file_on_branch(file_path, repo, default_branch, dataframe, assignees)
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
            update_breakdown = dataframe_update_breakdown(old_dataframe, dataframe)
            if any(item["has_update"] for item in update_breakdown.values()):
                print("Change detected")
            else:
                print("Nothing to update")
                pass
            branch = default_branch
            for key, value in update_breakdown.items():
                commit_message = key
                pr_message = ""
                if value["has_update"]:
                    print("Update: " + key)
                    if value["type"] == UpdateType.DRIFT and value["drift_context"]:
                        drift_evaluation = safe_drift_evaluator(
                            value["drift_context"], drift_evaluator
                        )
                        commit_message = key + "\n\n" + drift_evaluation["message"]
                        if drift_evaluation["should_alert"]:
                            checkout_branch_from_default_branch(repo, drift_branch)
                            pr_message = (
                                pr_message + "\n\n" + drift_evaluation["message"]
                            )
                            branch = drift_branch

                    update_file_with_retry(
                        repo=repo,
                        file_path=file_path,
                        commit_message=commit_message,
                        data=value["df"].to_csv(index=True, header=True),
                        branch=branch,
                    )
                    if pr_message != "":
                        create_pullrequest(
                            repo, branch, assignees, file_path, pr_message
                        )


def assert_file_exists(
    repo: Repository.Repository, file_path: str, ref: str
) -> Optional[ContentFile.ContentFile]:
    try:
        contents = repo.get_contents(file_path, ref=ref)
        assert not isinstance(contents, list), "pathfile returned multiple contents"
        return contents
    except GithubException as e:
        if e.status == 404:
            return None
        else:
            raise e


def update_file_with_retry(
    repo: Repository.Repository, file_path, commit_message, data, branch, max_retries=3
):
    retries = 0

    while retries < max_retries:
        try:
            content = assert_file_exists(repo, file_path, ref=branch)
            if content is None:
                response = repo.create_file(file_path, commit_message, data, branch)
                print(response["commit"].html_url)
            else:
                response = repo.update_file(
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


def create_file_on_branch(
    file_path: str,
    repo: Repository.Repository,
    branch: str,
    dataframe: pd.DataFrame,
    assignees: List[str],
):
    commit_message = "New data: " + file_path
    print("Commit: " + commit_message)
    repo.create_file(
        file_path, commit_message, dataframe.to_csv(index=True, header=True), branch
    )


def create_pullrequest(
    repo: Repository.Repository,
    branch: str,
    assignees: List[str],
    file_path: str,
    description_body: str,
):
    try:
        if len(assignees) > 0:
            pullrequest = repo.create_pull(
                title="New drift detected " + file_path,
                body=description_body,
                head=branch,
                base=repo.default_branch,
            )
            print("Pull request created: " + pullrequest.html_url)
            existing_assignees = assert_assignees_exists(repo, assignees)
            pullrequest.add_to_assignees(*existing_assignees)
        else:
            print("No assignees, skipping pull request creation")
    except GithubException as e:
        if e.status == 422:
            print("Pull request already exists, skipping...")
        else:
            raise e


def assert_branch_exist(repo: Repository.Repository, branch_name: str) -> None:
    try:
        branch = repo.get_branch(branch_name)
    except:
        branch = None

    # If the branch doesn't exist, create it
    if not branch:
        print(f"Branch {branch_name} doesn't exist, creating it...")
        create_git_branch(repo, branch_name)


def assert_assignees_exists(
    repo: Repository.Repository, assignees: List[str]
) -> List[str]:
    members = [collaborator.login for collaborator in repo.get_collaborators()]
    exising_assignees = []
    for assignee in assignees:
        if assignee not in members:
            print(f"Assignee {assignee} does not exist")
        else:
            exising_assignees.append(assignee)
    return exising_assignees


def get_valid_branch_name(filepath: str) -> str:
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


def create_git_branch(repo: Repository.Repository, branch_name: str):
    """
    Creates a new Git branch with the given name in the given repository.
    """
    # Get the default branch of the repository
    reported_branch = repo.get_branch(repo.default_branch)

    # Create a new reference to the default branch
    ref = repo.create_git_ref(f"refs/heads/{branch_name}", reported_branch.commit.sha)

    return ref


def find_date_column(df):
    date_columns = df.filter(like="date").columns
    if len(date_columns) > 0:
        return date_columns[0]
    else:
        return df.columns[0]


def checkout_branch_from_default_branch(repo: Repository.Repository, branch_name: str):
    assert_branch_exist(repo, branch_name)
    try:
        ref = repo.get_git_ref(f"heads/{branch_name}")
        ref.delete()
    except GithubException:
        pass

    """
    Checkout a branch from the default branch of the given repository.
    """
    # Get the default branch of the repository
    default_branch = repo.get_branch(repo.default_branch)
    print(
        "Checkout branch: " + branch_name, " from default branch:" + default_branch.name
    )

    # Create a new reference to the default branch

    try:
        ref = repo.create_git_ref(
            f"refs/heads/{branch_name}", default_branch.commit.sha
        )
    except GithubException:
        pass
    return


def get_monthly_file_path(file_path, month):
    directory, file_name = os.path.split(file_path)
    file_name, extension = os.path.splitext(file_name)

    new_file_name = f"{file_name}/{month}{extension}"

    new_file_path = os.path.join(directory, new_file_name)

    return new_file_path
