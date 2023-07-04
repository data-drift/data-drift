import time
from typing import Optional, List
import pandas as pd
from github import Github, Repository, ContentFile, GithubException
from datagit.dataset_helpers import compare_dataframes, sort_dataframe_on_first_column_and_assert_is_unique
import re
import base64
import io


def store_metric(ghClient: Github, dataframe: pd.DataFrame, filepath: str, assignees: List[str] = [], branch: Optional[str] = None, store_json: bool = True) -> None:
    print("Storing metric...")
    repo_orga, repo_name, file_path = filepath.split('/', 2)
    branch = branch or get_valid_branch_name(file_path)

    repo = ghClient.get_repo(repo_orga + "/" + repo_name)
    dataframe = sort_dataframe_on_first_column_and_assert_is_unique(dataframe)

    push_metric(dataframe, assignees, repo.default_branch,
                branch, store_json, file_path, repo)


def push_metric(dataframe, assignees, reported_branch, computed_branch, store_json, file_path, repo):
    contents = assert_file_exists(repo, file_path, ref=reported_branch)
    if contents is None:
        print("Metric not found, creating it on branch: " + reported_branch)
        create_file_on_branch(file_path, repo, reported_branch,
                              dataframe, assignees, store_json)
        pass
    else:
        print("Metric found, updating it on branch: " + reported_branch)
        date_column = find_date_column(dataframe)
        if contents.content is not None and date_column is not None:
            # Compare the contents of the file with the new contents and assert if it need 2 commits
            decoded_content = base64.b64decode(contents.content)
            content_string = decoded_content.decode('utf-8')
            old_dataframe = pd.read_csv(io.StringIO(content_string))
            try:
                old_dates = set(old_dataframe[date_column])
            except KeyError:
                old_dates = []
            new_dates = set(dataframe[date_column])
            already_stored_dates = new_dates.intersection(old_dates)
            new_dataframe = dataframe[~dataframe[date_column].isin(
                already_stored_dates)]
            old_data_with_freshdata = pd.concat(
                [old_dataframe, new_dataframe]).reset_index(drop=True)
            if len(new_dataframe) > 0:
                print("New data found")
                push_new_lines(
                    file_path, repo, reported_branch, old_data_with_freshdata, store_json)
            checkout_branch_from_default_branch(repo, computed_branch)

            if not old_data_with_freshdata.equals(dataframe.reset_index(drop=True)):
                print("Drift detected")

                push_drift_lines(file_path, repo, computed_branch,
                                 dataframe, store_json)
                print("Drift pushed")
                print("Creating pull request")
                description_body = f"Drift detected:\n" + \
                    compare_dataframes(old_data_with_freshdata,
                                       dataframe, "unique_key")
                create_pullrequest(repo, computed_branch,
                                   assignees, file_path, description_body)

    pass


def assert_file_exists(repo: Repository.Repository, file_path: str, ref: str) -> Optional[ContentFile.ContentFile]:
    try:
        contents = repo.get_contents(file_path, ref=ref)
        assert not isinstance(
            contents, list), "pathfile returned multiple contents"
        return contents
    except GithubException as e:
        if e.status == 404:
            return None
        else:
            raise e


def push_new_lines(file_path: str, repo: Repository.Repository, branch: str, dataframe: pd.DataFrame, store_json: bool):
    commit_message = "New data: " + file_path
    print("Commit: " + commit_message)
    update_file_with_retry(repo, file_path, commit_message,
                           dataframe.to_csv(index=False, header=True), branch)

    if store_json:
        json_file_path = file_path.replace(".csv", ".json")
        json_commit_message = "New data (json): " + json_file_path
        json_data = dataframe.to_json(
            orient="records", lines=True, date_format="iso")
        update_file_with_retry(repo, json_file_path,
                               json_commit_message, json_data, branch)


def push_drift_lines(file_path: str, repo: Repository.Repository, branch: str, dataframe: pd.DataFrame, store_json: bool):
    commit_message = "Drift: " + file_path
    print("Commit: " + commit_message)
    update_file_with_retry(repo, file_path, commit_message,
                           dataframe.to_csv(index=False, header=True), branch)

    if store_json:
        json_file_path = file_path.replace(".csv", ".json")
        json_commit_message = "Drift (json): " + json_file_path
        json_data = dataframe.to_json(
            orient="records", lines=True, date_format="iso")
        update_file_with_retry(repo, json_file_path,
                               json_commit_message, json_data, branch)


def update_file_with_retry(repo: Repository.Repository, file_path, commit_message, data, branch, max_retries=3):
    retries = 0

    while retries < max_retries:
        try:
            content = assert_file_exists(repo, file_path, ref=branch)
            if content is None:
                repo.create_file(file_path, commit_message, data, branch)
            else:
                repo.update_file(file_path, commit_message,
                                 data, content.sha, branch)
            return
        except GithubException as e:
            if e.status == 409:
                retries += 1
                time.sleep(1)  # Wait for 1 second before retrying
            else:
                raise e
    raise Exception(f"Failed to update file after {max_retries} retries")


def create_file_on_branch(file_path: str, repo: Repository.Repository, branch: str, dataframe: pd.DataFrame, assignees: List[str], store_json: bool):
    commit_message = "New data: " + file_path
    print("Commit: " + commit_message)
    repo.create_file(file_path, commit_message, dataframe.to_csv(
        index=False, header=True), branch)
    if store_json:
        json_file_path = file_path.replace(".csv", ".json")
        repo.create_file(json_file_path, "New data (json): " +
                         file_path, dataframe.to_json(orient="records", lines=True, date_format="iso"), branch)


def create_pullrequest(repo: Repository.Repository, branch: str, assignees: List[str], file_path: str, description_body: str):
    try:
        if len(assignees) > 0:
            pullrequest = repo.create_pull(
                title="New drift detected "+file_path, body=description_body, head=branch, base=repo.default_branch)
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


def assert_assignees_exists(repo: Repository.Repository, assignees: List[str]) -> List[str]:
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
    branch_name = re.sub(r'[^a-zA-Z0-9]+', '-', filepath)

    # Remove any leading or trailing hyphens
    branch_name = branch_name.strip('-')

    # Convert to lowercase
    branch_name = branch_name.lower()

    # Append a prefix
    branch_name = f'metric/{branch_name}'

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
    ref = repo.create_git_ref(
        f'refs/heads/{branch_name}', reported_branch.commit.sha)

    return ref


def find_date_column(df):
    date_columns = df.filter(like='date').columns
    if len(date_columns) > 0:
        return date_columns[0]
    else:
        return df.columns[0]


def checkout_branch_from_default_branch(repo: Repository.Repository, branch_name: str):
    assert_branch_exist(repo, branch_name)
    try:
        ref = repo.get_git_ref(f'heads/{branch_name}')
        ref.delete()
    except GithubException:
        pass

    """
    Checkout a branch from the default branch of the given repository.
    """
    # Get the default branch of the repository
    default_branch = repo.get_branch(repo.default_branch)

    # Create a new reference to the default branch
    ref = repo.create_git_ref(
        f'refs/heads/{branch_name}', default_branch.commit.sha)

    return ref
